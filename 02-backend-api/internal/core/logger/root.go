package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/constants"
)

func SetupLogger() {
	baseHandler := tint.NewTextHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.TimeOnly,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == constants.RequestIDKey {
				return slog.String("req_id", fmt.Sprintf("%s", a.Value.String()))
			}
			return a
		},
	})
	wrappedHandler := &ContextHandler{Handler: baseHandler}
	slog.SetDefault(slog.New(wrappedHandler))
}
