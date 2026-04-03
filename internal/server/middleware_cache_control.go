package server

import (
	"net/http"
)

type cacheControlWriter struct {
	http.ResponseWriter

	cacheControl string
}

func (w *cacheControlWriter) WriteHeader(statusCode int) {
	w.Header().Set("Cache-Control", w.cacheControl)
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *cacheControlWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func withCacheControl(value string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w = &cacheControlWriter{ResponseWriter: w, cacheControl: value}

		next.ServeHTTP(w, r)
	}
}
