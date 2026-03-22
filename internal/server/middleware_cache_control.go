package server

import (
	"net/http"
)

func withCacheControl(value string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", value)

			next.ServeHTTP(w, r)
		}
	}
}
