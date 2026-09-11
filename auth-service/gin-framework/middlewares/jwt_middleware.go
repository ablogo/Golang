package middlewares

import (
	"net/http"
	"strings"

	"src/services"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.Path, "auth") {
			ctx.Next()
			return
		}

		tokenHeader := ctx.GetHeader("Authorization")
		token := strings.Split(tokenHeader, " ")
		if len(token) <= 1 || strings.ToLower(token[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			ctx.Abort()
			return
		}

		userId, isValid := services.ValidateToken(token[1])
		if isValid {
			ctx.Set("user_id", userId)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
