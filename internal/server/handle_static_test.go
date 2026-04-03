package server

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	//go:embed static/packages.json
	packagesJSON string

	//go:embed static/robots.txt
	robotsTxt string
)

func TestStatic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		path            string
		wantContentType string
		wantBody        string
	}{
		{"packages_json", "/packages.json", "application/json", packagesJSON},
		{"robots_txt", "/robots.txt", "text/plain; charset=utf-8", robotsTxt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer()
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
			}

			if got := rec.Header().Get("Content-Type"); got != tt.wantContentType {
				t.Errorf("Content-Type header = %q, want %q", got, tt.wantContentType)
			}

			if got := rec.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}

			assertCacheControl(t, rec, defaultCacheControl)
		})

		t.Run(tt.name+"/conditional_get", func(t *testing.T) {
			t.Parallel()

			assertConditionalGet(t, tt.path, nil)
		})
	}
}
