package logger

import (
	"context"
	"log/slog"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/constants"
)

type ContextHandler struct {
	slog.Handler
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if reqID, ok := ctx.Value(constants.RequestIDKey).(string); ok {
			r.AddAttrs(slog.String("req_id", reqID))
		}
		// if userID, ok := ctx.Value(middleware.AuthUserID).(string); ok {
		// 	r.AddAttrs(slog.String("user_id", userID))
		// }
	}

	return h.Handler.Handle(ctx, r)
}
