package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestConditionalGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data map[string][]byte
		url  string
	}{
		{
			name: "packages",
			url:  "/packages.json",
		},
		{
			name: "robots",
			url:  "/robots.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modTime := time.Date(2016, 1, 2, 3, 4, 5, 6, time.UTC)

			srv := New(&stubStore{data: tt.data}, modTime)

			req1 := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
			rec1 := httptest.NewRecorder()

			srv.ServeHTTP(rec1, req1)

			if rec1.Code != http.StatusOK {
				t.Errorf("first status = %d, want %d", rec1.Code, http.StatusOK)
			}

			gotLM := rec1.Header().Get("Last-Modified")

			wantLM := modTime.Truncate(time.Second).Format(http.TimeFormat)
			if gotLM != wantLM {
				t.Fatalf("Last-Modified header = %q, want %q", gotLM, wantLM)
			}

			rec1CC := rec1.Header().Get("Cache-Control")

			b1, err := io.ReadAll(rec1.Body)
			if err != nil {
				t.Fatalf("io.ReadAll(rec1.Body) unexpected error: %v", err)
			}
			if len(b1) == 0 {
				t.Error("first body len = 0, want > 0")
			}

			req2 := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
			req2.Header.Set("If-Modified-Since", gotLM)
			rec2 := httptest.NewRecorder()

			srv.ServeHTTP(rec2, req2)

			if rec2.Code != http.StatusNotModified {
				t.Errorf("second status = %d, want %d", rec2.Code, http.StatusNotModified)
			}

			if got := rec2.Header().Get("Cache-Control"); got != rec1CC {
				t.Errorf("second Cache-Control header = %q, want %q", got, rec1CC)
			}

			hs := []string{
				"Content-Type",
				"Content-Length",
				"Content-Encoding",
			}
			for _, h := range hs {
				if got := rec2.Header().Get(h); got != "" {
					t.Errorf("second %q header = %q, want empty", h, got)
				}
			}

			b2, err := io.ReadAll(rec2.Body)
			if err != nil {
				t.Fatalf("io.ReadAll(rec2.Body) unexpected error: %v", err)
			}
			if len(b2) != 0 {
				t.Errorf("second body = %q, want empty", b2)
			}
		})
	}
}

func TestConditionalGet_NoLastModifiedOnError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		url  string
	}{
		{
			name: "not-found",
			url:  "/not-found-foo-bar",
		},
		{
			name: "p2/stable/not-found",
			url:  "/p2/foo/bar.json",
		},
		{
			name: "p2/dev/not-found",
			url:  "/p2/foo/bar~dev.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			srv := newTestServer()
			req := httptest.NewRequest(http.MethodGet, tt.url, http.NoBody)
			rec := httptest.NewRecorder()

			srv.ServeHTTP(rec, req)

			if rec.Code < 400 {
				t.Errorf("status = %d, want >= 400", rec.Code)
			}

			if got := rec.Header().Get("Last-Modified"); got != "" {
				t.Errorf("Last-Modified header = %q, want empty", got)
			}
		})
	}
}
