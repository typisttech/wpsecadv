package server

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"
)

func TestUpRoute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
	}{
		{"get", http.MethodGet},
		{"post", http.MethodPost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			synctest.Test(t, func(t *testing.T) {
				dst := time.Date(2029, 12, 13, 14, 15, 16, 17, time.UTC)

				time.Sleep(time.Until(dst))

				srv := newTestServer()
				req := httptest.NewRequest(tt.method, "/up", http.NoBody)
				rec := httptest.NewRecorder()

				srv.ServeHTTP(rec, req)

				if rec.Code != http.StatusOK {
					t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
				}

				if lm := rec.Header().Get("Last-Modified"); lm != "" {
					t.Errorf("Last-Modified header = %q, want empty", lm)
				}

				gotCC := rec.Header().Get("Cache-Control")
				wantCC := "no-store"
				if gotCC != wantCC {
					t.Errorf("Cache-Control header = %q, want %q", gotCC, wantCC)
				}

				ct := rec.Header().Get("Content-Type")
				if ct != "application/json" {
					t.Errorf("Content-Type header = %q, want %q", ct, "application/json")
				}

				var got struct {
					Status    string `json:"status"`
					Timestamp string `json:"timestamp"`
				}

				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("json.Unmarshal() unexpected error: %v", err)
				}

				if got.Status != "up" {
					t.Errorf("status = %q, want %q", got.Status, "up")
				}

				gotT, err := time.Parse(time.RFC3339, got.Timestamp)
				if err != nil {
					t.Fatalf("time.Parse() unexpected error: %v", err)
				}

				wantT := dst.Truncate(time.Second)
				if !gotT.Equal(wantT) {
					t.Errorf("timestamp = %v, want %v", gotT, wantT)
				}
			})
		})
	}
}
