package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kalpeshWani222/olx-api/internal/httpx"
	"github.com/kalpeshWani222/olx-api/internal/token"
)


const (
	userIDKey ctxKey = iota + 1
)

//auth middleware
func RequireAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpx.Error(w, http.StatusUnauthorized, "missing authorization header", httpx.CodeUnauthenticated)
				return
			}

			//Check that it starts with "Bearer "
			tokenStr, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found {
				httpx.Error(w, http.StatusUnauthorized, "invalid authorization format, use: Bearer <token>", httpx.CodeUnauthenticated)
				return
			}

			// Parse and validate the JWT
			claims, err := token.Parse(tokenStr, secret)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "invalid or expired token", httpx.CodeUnauthenticated)
				return
			}

			// Store the userID in the request context
			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext reads the authenticated user's ID from the context.
func UserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}
