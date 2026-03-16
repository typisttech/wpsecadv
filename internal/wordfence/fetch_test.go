package wordfence

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/typisttech/wpsecadv/internal"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any) {}
func (noopLogger) Warn(string, ...any) {}

func sortRecords() cmp.Option {
	return cmpopts.SortSlices(func(a, b internal.Record) bool {
		if a.Slug != b.Slug {
			return a.Slug < b.Slug
		}
		return a.Kind < b.Kind
	})
}

func TestFetch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		body        string
		wantRecords []internal.Record
	}{
		{
			name: "single_vulnerability_produces_records",
			body: `{
				"abc-123": {
					"id": "abc-123",
					"software": [{
						"type": "plugin",
						"slug": "my-plugin",
						"affected_versions": {
							"*-1.0": {"from_version": "*", "from_inclusive": true, "to_version": "1.0", "to_inclusive": true}
						}
					}]
				}
			}`,
			wantRecords: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "my-plugin",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/abc-123/my-plugin",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "abc-123"}},
						AffectedVersions: "<=1.0",
					},
				},
			},
		},
		{
			name: "multiple_vulnerabilities_produce_records",
			body: `{
				"aaa-111": {
					"id": "aaa-111",
					"software": [{
						"type": "plugin",
						"slug": "plugin-a",
						"affected_versions": {
							"*-1.0": {"from_version": "*", "from_inclusive": true, "to_version": "1.0", "to_inclusive": true}
						}
					}]
				},
				"bbb-222": {
					"id": "bbb-222",
					"software": [{
						"type": "theme",
						"slug": "theme-b",
						"affected_versions": {
							"*-2.0": {"from_version": "*", "from_inclusive": true, "to_version": "2.0", "to_inclusive": true}
						}
					}]
				}
			}`,
			wantRecords: []internal.Record{
				{
					Kind: internal.KindPlugin,
					Slug: "plugin-a",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/aaa-111/plugin-a",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "aaa-111"}},
						AffectedVersions: "<=1.0",
					},
				},
				{
					Kind: internal.KindTheme,
					Slug: "theme-b",
					Advisory: internal.Advisory{
						ID:               "WPSECADV/WF/bbb-222/theme-b",
						Sources:          []internal.Source{{Name: "Wordfence", ID: "bbb-222"}},
						AffectedVersions: "<=2.0",
					},
				},
			},
		},
		{
			name:        "vulnerability_with_no_valid_software_produces_no_records",
			body:        `{"abc-123": {"id": "abc-123", "software": []}}`,
			wantRecords: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tt.body))
			}))
			t.Cleanup(srv.Close)

			rc, errc := Fetch(t.Context(), noopLogger{}, srv.Client(), srv.URL, "")

			var got []internal.Record
			for r := range rc {
				got = append(got, r)
			}
			err := <-errc
			if err != nil {
				t.Fatalf("Fetch() unexpected error: %v", err)
			}

			if diff := cmp.Diff(tt.wantRecords, got, sortRecords(), cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("Fetch() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFetch_Non200Response(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	_, errc := Fetch(t.Context(), noopLogger{}, srv.Client(), srv.URL, "")

	err := <-errc
	if err == nil {
		t.Error("Fetch() unexpected success")
	}
}

func TestFetch_ContextCancelled_BeforeFetch(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"abc-123": {"id": "abc-123"}}`))
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	_, errc := Fetch(ctx, noopLogger{}, srv.Client(), srv.URL, "")

	err := <-errc

	if err == nil {
		t.Error("Fetch() unexpected success")
	}
}

func TestFetch_ContextCancelled_MidStream(t *testing.T) {
	t.Parallel()

	body := `{
		"id-a": {"id": "id-a", "software": [{"type": "plugin", "slug": "plugin-a", "affected_versions": {"*-1.0": {"from_version": "*", "from_inclusive": true, "to_version": "1.0", "to_inclusive": true}}}]},
		"id-b": {"id": "id-b", "software": [{"type": "plugin", "slug": "plugin-b", "affected_versions": {"*-2.0": {"from_version": "*", "from_inclusive": true, "to_version": "2.0", "to_inclusive": true}}}]},
		"id-c": {"id": "id-c", "software": [{"type": "plugin", "slug": "plugin-c", "affected_versions": {"*-3.0": {"from_version": "*", "from_inclusive": true, "to_version": "3.0", "to_inclusive": true}}}]}
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	rc, errc := Fetch(ctx, noopLogger{}, srv.Client(), srv.URL, "")

	// Read one record to confirm the stream has started, then cancel and stop
	// reading. The goroutine will be blocked on rc <- r and pick ctx.Done().
	<-rc
	cancel()

	err := <-errc
	if err == nil {
		t.Error("Fetch() unexpected success")
	}
}
