package server

import (
	"net/http"
	"time"
)

func handleUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	ts := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	w.Write([]byte(`{"status":"up","timestamp":"` + ts + `"}`))
}
