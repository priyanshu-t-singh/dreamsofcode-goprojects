package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/04-url-shortener/internal/constants"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the incoming request already has an X-Request-ID header
		// e.g., from an API gateway/load balancer
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateID()
		}

		w.Header().Set("X-Request-ID", reqID)

		// Create context containing the Request ID
		ctx := context.WithValue(r.Context(), constants.RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(r *http.Request) string {
	if reqID, ok := r.Context().Value(constants.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func generateID() string {
	return uuid.New().String()
}
