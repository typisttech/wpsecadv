package server

import (
	"embed"
	"io/fs"
	"net/http"
	"time"
)

//go:embed  static/packages.json static/robots.txt
var static embed.FS

func addRoutes(mux *http.ServeMux, store AdvisoriesMarshaler, modTime time.Time) {
	m := withConditionalGet(modTime)

	hAdvs := handleAdvisories(store)
	mux.HandleFunc("GET /api/security-advisories/{$}", withCacheControl("max-age=3600")(m(hAdvs)))
	mux.HandleFunc("POST /api/security-advisories/{$}", withCacheControl("max-age=3600")(hAdvs))

	// Health check.
	hUp := withCacheControl("no-store")(http.HandlerFunc(handleUp))
	mux.HandleFunc("GET /up", hUp)
	mux.HandleFunc("POST /up", hUp)

	// In case someone clicks form composer.json, redirect them to the GitHub repo.
	mux.Handle("GET /{$}", http.RedirectHandler("https://github.com/typisttech/wpsecadv", http.StatusFound))

	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	hFS := m(http.FileServerFS(sub))
	mux.Handle("GET /packages.json", hFS)
	mux.Handle("GET /robots.txt", hFS)
}
