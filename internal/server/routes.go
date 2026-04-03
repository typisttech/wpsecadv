package server

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed  static/packages.json static/robots.txt
var static embed.FS

func addRoutes(mux *http.ServeMux, store AdvisoriesMarshaler) {
	mux.HandleFunc("GET /p2/{vendor}/{file}", handleP2(store))

	sub, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	hFS := http.FileServerFS(sub)
	mux.Handle("GET /packages.json", hFS)
	mux.Handle("GET /robots.txt", hFS)

	// In case someone clicks form composer.json, redirect them to the GitHub repo.
	mux.Handle("GET /{$}", http.RedirectHandler("https://github.com/typisttech/wpsecadv", http.StatusFound))
}
