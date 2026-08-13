package web

import (
	"net/http"
)

func InitiateFileServer(router *http.ServeMux) {
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	router.Handle("/static/", http.FileServer(http.FS(staticFiles)))
}
