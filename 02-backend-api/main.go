package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/arithmetics"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/constants"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/database"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/logger"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/middleware"
)

func main() {
	logger.SetupLogger()
	router := http.NewServeMux()
	v1 := http.NewServeMux()

	db, err := database.Open(constants.SqliteDatabasePath)
	if err != nil {
		slog.Error("db failed to open", "err", err)
		return
	}

	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Calculator App is running")
	})

	v1.Handle("/", middleware.IsAuthenticated(arithmetics.GetRoutes(db)))
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	server := http.Server{
		Addr: ":8080",
		Handler: middleware.CreateStack(
			middleware.RateLimiter,
			middleware.RequestID,
			middleware.AllowCors,
			middleware.Logging,
		)(router),
	}

	go func() {
		slog.Info("Server listening on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to serve", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
	if err := db.Close(); err != nil {
		slog.Error("error closing database", "err", err)
	}
	slog.Info("shutdown complete")
}
