package middleware

import (
	"context"
	"net/http"
	"strings"

	"src/services"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isLogin := strings.Contains(r.RequestURI, "auth")
		if isLogin {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")

		token := strings.Split(tokenHeader, " ")
		if len(token) <= 1 || strings.ToLower(token[0]) != "bearer" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		userId, isValid := services.ValidateToken(token[1])
		if isValid {
			ctx := context.WithValue(r.Context(), "userId", userId)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

	})

}
