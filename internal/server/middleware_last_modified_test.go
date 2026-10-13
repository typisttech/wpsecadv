package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func assertLastModified(t *testing.T, rec *httptest.ResponseRecorder, want time.Time) {
	t.Helper()

	gotLM := rec.Header().Get("Last-Modified")
	wantLM := want.Format(http.TimeFormat)

	if gotLM != wantLM {
		t.Fatalf("Last-Modified header = %q, want %q", gotLM, wantLM)
	}
}
