package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/constants"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/core/database"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/core/logger"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/middleware"
	us "github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/url"
)

type PageData struct {
	BaseURL     string
	ShortURLId  string
	OriginalURL string
}

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

	repo := us.NewRepository(db)
	urlHandler := us.NewHandler(repo)
	router.HandleFunc("POST /api/shorten", urlHandler.SaveShortURL)
	router.HandleFunc("GET /{pageId}", urlHandler.RedirectToOriginalURL)

	router.HandleFunc("GET /r/{pageId}", func(w http.ResponseWriter, r *http.Request) {
		pageId := r.PathValue("pageId")

		originalURL, err := repo.GetOriginalURL(pageId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "result.html", PageData{
			BaseURL:     constants.BaseURL,
			OriginalURL: originalURL,
			ShortURLId:  pageId,
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
