package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/arithmetics"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/database"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/logger"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/middleware"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/web"
)

func main() {
	logger.SetupLogger()
	router := http.NewServeMux()
	v1 := http.NewServeMux()
	web.InitiateFileServer(router)

	db, err := database.SetupDatabase()
	if err != nil {
		slog.Error("startup failed", "err", err)
		os.Exit(1)
	}

	v1.Handle("/", middleware.IsAuthenticated(arithmetics.GetRoutes(db)))
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

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
