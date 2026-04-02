package server

import (
	"net/http"
	"strings"
)

func handleP2(store AdvisoriesMarshaler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vendor := r.PathValue("vendor")
		vendor = strings.ToLower(vendor)

		file := r.PathValue("file")
		file = strings.ToLower(file)
		if !strings.HasSuffix(file, ".json") {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(file, "~dev.json") {
			http.NotFound(w, r)
			return
		}

		slug := strings.TrimSuffix(file, ".json")

		advisories, err := store.MarshalAdvisoriesFor(vendor, slug)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"packages":[],"security-advisories":`))
		//gosec:disable G705 -- Advisories bytes originate from trusted embedded asset files
		w.Write(advisories)
		w.Write([]byte(`}`))
	}
}
