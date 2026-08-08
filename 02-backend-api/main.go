package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/arithmetics"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	handler := &arithmetics.Handler{}
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Calculator App is running")
	})

	router.HandleFunc("POST /add", handler.Add)
	router.HandleFunc("POST /subtract", handler.Subtract)
	router.HandleFunc("POST /multiply", handler.Multiply)
	router.HandleFunc("POST /divide", handler.Divide)
	router.HandleFunc("POST /sum", handler.Sum)

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.AllowCors(middleware.Logging(logger, router)),
	}

	logger.Info("Server listening on port :8080")
	server.ListenAndServe()
}
