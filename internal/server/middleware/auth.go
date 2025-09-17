// Package middleware предоставляет функционал для обработчиков middleware.
package middleware

import (
	"context"
	"net/http"

	"github.com/mr-filatik/go-goph-keeper/internal/server/crypto/jwt"
)

type ctxKey string

const userIDKey ctxKey = "uid"

// WithUserID добавляет в контекст идентификатор пользователя.
func WithUserID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, userIDKey, uid)
}

// GetUserID получает идентификатор пользователя из контекста.
func GetUserID(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(userIDKey).(string)

	return value, ok
}

// RequireAuth представляет middleware для авторизации пользователей.
func RequireAuth(enc *jwt.Encryptor, next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		token, err := enc.ValidateTokenBearer(req.Header.Get("Authorization"))
		if err != nil {
			http.Error(resp, "invalid token", http.StatusUnauthorized)

			return
		}

		uid, err := enc.GetClaimUserIDFromToken(token)
		if err != nil {
			http.Error(resp, "invalid claims", http.StatusUnauthorized)

			return
		}

		next(resp, req.WithContext(WithUserID(req.Context(), uid)))
	}
}
