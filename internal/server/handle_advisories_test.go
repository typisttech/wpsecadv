package server

import (
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/typisttech/wpsecadv/internal"
)

type body struct {
	Status     string                         `json:"status"`
	Message    string                         `json:"message"`
	Advisories map[string][]internal.Advisory `json:"advisories"`
}

func TestAdvisories(t *testing.T) {
	t.Parallel()

	const path = "/api/security-advisories/"

	tests := []struct {
		name         string
		data         map[string][]byte
		packages     []string
		updatedSince string
		wantCode     int
		want         body
	}{
		{
			name:         "missing_packages",
			packages:     nil,
			updatedSince: "",
			wantCode:     http.StatusBadRequest,
			want: body{
				Status:  "error",
				Message: `Missing array of package names as the "packages" parameter`,
			},
		},
		{
			name:         "empty_packages",
			packages:     []string{""},
			updatedSince: "",
			wantCode:     http.StatusBadRequest,
			want: body{
				Status:  "error",
				Message: `Missing array of package names as the "packages" parameter`,
			},
		},
		{
			name:         "updated_since_only",
			packages:     nil,
			updatedSince: "2026-01-01",
			wantCode:     http.StatusBadRequest,
			want: body{
				Status:  "error",
				Message: `The "updatedSince" parameter is not supported.`,
			},
		},
		{
			name:         "updated_since_and_packages",
			packages:     []string{"any-vendor/woocommerce"},
			updatedSince: "2026-01-01",
			wantCode:     http.StatusBadRequest,
			want: body{
				Status:  "error",
				Message: `Pass only one of "updatedSince" OR "packages" parameters, they cannot be provided together.`,
			},
		},
		{
			name: "single_valid_package",
			data: map[string][]byte{
				"any-vendor/woocommerce": []byte(`[{"advisoryId":"WPSECADV/1"}]`),
			},
			packages:     []string{"any-vendor/woocommerce"},
			updatedSince: "",
			wantCode:     http.StatusOK,
			want: body{
				Advisories: map[string][]internal.Advisory{
					"any-vendor/woocommerce": {{ID: "WPSECADV/1"}},
				},
			},
		},
		{
			name:         "skip_invalid_package",
			packages:     []string{"invalid-package"},
			updatedSince: "",
			wantCode:     http.StatusOK,
			want:         body{Advisories: map[string][]internal.Advisory{}},
		},
		{
			name:         "skip_not_found_package",
			packages:     []string{"wpackagist-plugin/non-existent-plugin"},
			updatedSince: "",
			wantCode:     http.StatusOK,
			want:         body{Advisories: map[string][]internal.Advisory{}},
		},
		{
			name: "mixed_valid_invalid",
			data: map[string][]byte{
				"any-vendor/woocommerce": []byte(`[{"advisoryId":"WPSECADV/1"}]`),
			},
			packages: []string{
				"invalid",
				"wpackagist-plugin/non-existent-plugin",
				"any-vendor/woocommerce",
			},
			updatedSince: "",
			wantCode:     http.StatusOK,
			want: body{
				Advisories: map[string][]internal.Advisory{
					"any-vendor/woocommerce": {{ID: "WPSECADV/1"}},
				},
			},
		},
		{
			name: "mixed_multiple_valid_invalid",
			data: map[string][]byte{
				"any-vendor/woocommerce":           []byte(`[{"advisoryId":"WPSECADV/1"}]`),
				"foo/bar":                          []byte(`[{"advisoryId":"WPSECADV/2"}]`),
				"wpackagist-plugin/example-plugin": []byte(`[{"advisoryId":"WPSECADV/3"}]`),
			},
			packages: []string{
				"any-vendor/woocommerce",
				"wpackagist-plugin/non-existent-plugin",
				"wpackagist-plugin/example-plugin",
			},
			updatedSince: "",
			wantCode:     http.StatusOK,
			want: body{
				Advisories: map[string][]internal.Advisory{
					"any-vendor/woocommerce":           {{ID: "WPSECADV/1"}},
					"wpackagist-plugin/example-plugin": {{ID: "WPSECADV/3"}},
				},
			},
		},
	}

	for _, tt := range tests {
		sq := url.Values{}
		for _, p := range tt.packages {
			sq.Add("packages[]", p)
		}

		num := url.Values{}
		for i, p := range tt.packages {
			num.Set(fmt.Sprintf("packages[%d]", i), p)
		}

		if tt.updatedSince != "" {
			sq.Set("updatedSince", tt.updatedSince)
			num.Set("updatedSince", tt.updatedSince)
		}

		sqEnc := sq.Encode()
		numEnc := num.Encode()

		t.Run("get/squared/"+tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, path+"?"+sqEnc, http.NoBody)
			assertAdvisoriesResponse(t, tt.data, req, tt.wantCode, tt.want)
		})

		t.Run("get/numbered/"+tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, path+"?"+numEnc, http.NoBody)
			assertAdvisoriesResponse(t, tt.data, req, tt.wantCode, tt.want)
		})

		t.Run("post/squared/"+tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(sqEnc))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			assertAdvisoriesResponse(t, tt.data, req, tt.wantCode, tt.want)
		})

		t.Run("post/numbered/"+tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(numEnc))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			assertAdvisoriesResponse(t, tt.data, req, tt.wantCode, tt.want)
		})
	}
}

func assertAdvisoriesResponse(t *testing.T, data map[string][]byte, req *http.Request, wantCode int, want body) {
	t.Helper()

	rec := httptest.NewRecorder()
	srv := newTestServerWithData(data)

	srv.ServeHTTP(rec, req)

	if rec.Code != wantCode {
		t.Errorf("status = %d, want %d", rec.Code, wantCode)
	}

	var got body
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal() unexpected error: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("response mismatch (-want +got):\n%s", diff)
	}
}
