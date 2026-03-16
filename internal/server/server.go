package server

import (
	"net/http"
	"time"
)

func New(store AdvisoriesMarshaler, modTime time.Time) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, store, modTime)

	return mux
}
