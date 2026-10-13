package server

import (
	"net/http"
	"time"
)

func withLastModified(modTime time.Time, next http.Handler) http.HandlerFunc {
	lm := modTime.UTC().Format(http.TimeFormat)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Last-Modified", lm)
		next.ServeHTTP(w, r)
	}
}
