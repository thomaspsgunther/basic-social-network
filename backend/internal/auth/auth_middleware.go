package auth

import (
	"context"
	"fmt"
	"net/http"

	"y_net/internal/logger"
	"y_net/internal/services/users"
	"y_net/pkg/jwt"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	username string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Validate jwt token
			tokenStr := header
			id, err := jwt.ParseToken(tokenStr)
			if err != nil {
				err := fmt.Errorf("invalid token")

				logger.ServerLogger.Warn(err.Error())

				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Create user and check if user exists in db
			user := users.User{ID: id}
			username, err := users.GetUsernameByUserID(id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.Username = username
			// Put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			// And call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *users.User {
	raw, _ := ctx.Value(userCtxKey).(*users.User)
	return raw
}
