package server

import (
	"net/http"
	"time"
)

func New(store AdvisoriesMarshaler, modTime time.Time) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)
	var handler http.Handler = mux
	handler = withLastModified(modTime, handler)
	handler = withCacheControl(defaultCacheControl, handler)

	final := http.NewServeMux()
	final.Handle("/", handler)

	final.HandleFunc("GET /up", withCacheControl("no-store", http.HandlerFunc(handleUp)))

	return final
}
