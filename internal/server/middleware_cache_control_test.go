package server

import (
	"net/http/httptest"
	"testing"
)

func assertCacheControl(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := rec.Header().Get("Cache-Control"); got != want {
		t.Errorf("Cache-Control header = %q, want %q", got, want)
	}
}
