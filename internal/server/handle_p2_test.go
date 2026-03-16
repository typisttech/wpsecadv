package server

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestP2(t *testing.T) {
	t.Parallel()

	methods := []string{http.MethodGet, http.MethodPost}

	tests := []struct {
		name string
		path string
	}{
		{"stable", "/p2/foo/bar.json"},
		{"dev", "/p2/foo/bar~dev.json"},
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					t.Parallel()

					srv := newTestServer()
					req := httptest.NewRequest(method, tt.path, http.NoBody)
					rec := httptest.NewRecorder()

					srv.ServeHTTP(rec, req)

					if rec.Code != http.StatusNotFound {
						t.Errorf("status code = %d, want %d", rec.Code, http.StatusNotFound)
					}
				})
			}
		})
	}
}
