package server

import (
	"net/http"
	"time"
)

const defaultCacheControl = "max-age=86400"

func New(store AdvisoriesMarshaler, modTime time.Time) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store)
	var handler http.Handler = mux
	handler = withConditionalGet(modTime, handler)
	handler = withCacheControl(defaultCacheControl, handler)

	final := http.NewServeMux()
	final.Handle("/", handler)

	final.HandleFunc("GET /up", withCacheControl("no-store", http.HandlerFunc(handleUp)))
	final.HandleFunc("POST /up", withCacheControl("no-store", http.HandlerFunc(handleUp)))

	return final
}
