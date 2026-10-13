package server

import (
	"net/http"
)

const defaultCacheControl = "public, max-age=43200, stale-while-revalidate=3600"

func withCacheControl(value string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", value)
		next.ServeHTTP(w, r)
	}
}
