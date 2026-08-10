package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/02-backend-api/internal/core/model"
)

const AuthUserID = "middleware.auth.userID"

func writeUnauthenticated(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.ErrorResponse{Error: msg})
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeUnauthenticated(w, "Unauthorized: Missing Authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeUnauthenticated(w, "Unauthorized: Invalid authorization format")
			return
		}

		token := parts[1]
		userID, err := validateToken(token)
		if err != nil {
			writeUnauthenticated(w, "Unauthorized: Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserID, userID)
		UpdateContext(r, ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(token string) (string, error) {
	if token == "c2VjcmV0LXRva2VuLTEyMzQ=" {
		return "user_1234", nil
	}
	return "", fmt.Errorf("invalid token")
}
