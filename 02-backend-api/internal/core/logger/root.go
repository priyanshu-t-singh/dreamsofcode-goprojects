package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func SetupLogger() {
	baseHandler := tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
	})
	wrappedHandler := &ContextHandler{Handler: baseHandler}
	slog.SetDefault(slog.New(wrappedHandler))
}
