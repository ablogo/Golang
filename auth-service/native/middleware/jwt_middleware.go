package middleware

import (
	"context"
	"net/http"
	"strings"

	"src/services"
)

func JWTMiddleware(next http.Handler) http.Handler {
	type userKey string
	const userIdKey userKey = "userId"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")

		token := strings.Split(tokenHeader, " ")
		if len(token) <= 1 || token[1] == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		userId, is_valid := services.ValidateToken(token[1])
		if is_valid {
			ctx := context.WithValue(r.Context(), userIdKey, userId)
			r.WithContext(ctx)

			next.ServeHTTP(w, r)

		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	})

}
