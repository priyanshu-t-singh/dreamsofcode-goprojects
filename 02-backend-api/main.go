package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/arithmetics"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/logger"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/middleware"
)

func main() {
	logger.SetupLogger()
	router := http.NewServeMux()
	v1 := http.NewServeMux()

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Calculator App is running")
	})

	v1.Handle("/", arithmetics.GetRoutes())
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	server := http.Server{
		Addr: ":8080",
		Handler: middleware.CreateStack(
			middleware.RequestID,
			middleware.AllowCors,
			middleware.Logging,
		)(router),
	}

	slog.Info("Server listening on port :8080")
	server.ListenAndServe()
}
