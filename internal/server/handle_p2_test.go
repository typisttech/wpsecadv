package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestP2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     map[string][]byte
		path     string
		wantCode int
		wantCT   string
		wantBody string
	}{
		{
			name: "known/stable",
			data: map[string][]byte{
				"foo/bar": []byte(`[{"advisoryId":"WPSECADV/1","affectedVersions":">=1.0.0"}]`),
			},
			path:     "/p2/foo/bar.json",
			wantCode: http.StatusOK,
			wantCT:   "application/json",
			wantBody: `{"packages":[],"security-advisories":[{"advisoryId":"WPSECADV/1","affectedVersions":">=1.0.0"}]}`,
		},
		{
			name: "known/dev",
			data: map[string][]byte{
				"foo/bar": []byte(`[{"advisoryId":"WPSECADV/1","affectedVersions":">=1.0.0"}]`),
			},
			path:     "/p2/foo/bar~dev.json",
			wantCode: http.StatusNotFound,
		},
		{
			name: "unknown/stable",
			data: map[string][]byte{
				"foo/bar": []byte(`[{"advisoryId":"WPSECADV/1","affectedVersions":">=1.0.0"}]`),
			},
			path:     "/p2/baz/qux.json",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "unknown/dev",
			path:     "/p2/foo/bar~dev.json",
			wantCode: http.StatusNotFound,
		},
		{name: "no_vendor/stable", path: "/p2/bar.json", wantCode: http.StatusNotFound},
		{name: "no_vendor/dev", path: "/p2/bar~dev.json", wantCode: http.StatusNotFound},
		{name: "no_file", path: "/p2/foo/", wantCode: http.StatusNotFound},
		{name: "no_slug/stable", path: "/p2/foo/.json", wantCode: http.StatusNotFound},
		{name: "no_slug/dev", path: "/p2/foo/~dev.json", wantCode: http.StatusNotFound},
		{name: "no_extension/stable", path: "/p2/foo/bar", wantCode: http.StatusNotFound},
		{name: "no_extension/dev", path: "/p2/foo/bar~dev", wantCode: http.StatusNotFound},
		{name: "not_json/stable", path: "/p2/foo/bar.txt", wantCode: http.StatusNotFound},
		{name: "not_json/dev", path: "/p2/foo/bar~dev.txt", wantCode: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServerWithData(tt.data)
			req := httptest.NewRequest(http.MethodGet, tt.path, http.NoBody)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)

			if rec.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantCode)
			}

			gotCT := rec.Header().Get("Content-Type")
			if tt.wantCT != "" && gotCT != tt.wantCT {
				t.Errorf("Content-Type = %q, want %q", gotCT, tt.wantCT)
			}

			gotBody := rec.Body.String()
			if tt.wantBody != "" && gotBody != tt.wantBody {
				t.Errorf("body = %q, want %q", gotBody, tt.wantBody)
			}
		})
	}
}
