package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed index.html css js
var staticFS embed.FS

func RegisterRoutes(mux *http.ServeMux) {
	cssFS, err := fs.Sub(staticFS, "css")
	if err != nil {
		panic("failed to load embedded css: " + err.Error())
	}
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.FS(cssFS))))

	jsFS, err := fs.Sub(staticFS, "js")
	if err != nil {
		panic("failed to load embedded js: " + err.Error())
	}
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.FS(jsFS))))

	mux.HandleFunc("/", serveIndex)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data, err := staticFS.ReadFile("index.html")
	if err != nil {
		http.Error(w, "failed to load page", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}
