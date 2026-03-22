package server

import (
	"net/http"
	"time"
)

func New(store AdvisoriesMarshaler, modTime time.Time) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store, modTime)

	var handler http.Handler = mux
	handler = withCacheControl("max-age=86400")(handler)

	return handler
}
