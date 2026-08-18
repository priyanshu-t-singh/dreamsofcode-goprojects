package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/core/database"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/core/logger"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/middleware"
)

type PageData struct {
	BaseURL    string
	ShortURLId string
}

var baseURL string = "http://localhost:8080"

func main() {
	logger.SetupLogger()

	db, err := database.SetupDatabase()
	if err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}

	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", PageData{})
	})

	router.HandleFunc("GET /r/{pageId}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("pageId")
		tmpl.ExecuteTemplate(w, "result.html", PageData{
			BaseURL:    baseURL,
			ShortURLId: id,
		})
	})

	server := &http.Server{
		Addr: ":8080",
		Handler: middleware.CreateStack(
			middleware.RateLimiter,
			middleware.RequestID,
			middleware.AllowCors,
			middleware.Logging,
		)(router),
	}

	database.RunServerWithGracefulShutdown(server, db)
}
