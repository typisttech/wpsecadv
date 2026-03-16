package server

import (
	"net/http"
	"time"
)

type lastModifiedWriter struct {
	http.ResponseWriter

	lastModified string
}

func (l *lastModifiedWriter) WriteHeader(statusCode int) {
	if statusCode == http.StatusOK {
		l.Header().Set("Last-Modified", l.lastModified)
	}

	l.ResponseWriter.WriteHeader(statusCode)
}

func withConditionalGet(modTime time.Time) func(http.Handler) http.HandlerFunc {
	ours := modTime.UTC().Format(http.TimeFormat)

	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)

				return
			}

			w = &lastModifiedWriter{ResponseWriter: w, lastModified: ours}

			theirs := r.Header.Get("If-Modified-Since")
			if theirs == ours {
				writeNotModified(w)

				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

func writeNotModified(w http.ResponseWriter) {
	// RFC 7232 section 4.1
	h := w.Header()
	delete(h, "Content-Type")
	delete(h, "Content-Length")
	delete(h, "Content-Encoding")
	if h.Get("Etag") != "" {
		delete(h, "Last-Modified")
	}
	w.WriteHeader(http.StatusNotModified)
}
