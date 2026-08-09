package logger

import (
	"context"
	"log/slog"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/constants"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if reqID, ok := ctx.Value(constants.RequestIDKey).(string); ok {
			r.AddAttrs(slog.String(constants.RequestIDKey, reqID))
		}
	}

	return h.Handler.Handle(ctx, r)
}
