package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/constants"
)

func SetupDatabase() (*sql.DB, error) {
	db, err := Open(constants.SqliteDatabasePath)
	if err != nil {
		return nil, fmt.Errorf("db failed to open: %w", err)
	}
	return db, nil
}

// RunServerWithGracefulShutdown starts the server, blocks until SIGINT/SIGTERM,
// then shuts down the server and closes the DB together.
func RunServerWithGracefulShutdown(server *http.Server, db *sql.DB) {
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
