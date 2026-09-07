package middlewares

import (
	"net/http"
	"strings"

	"src/services"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		/*if strings.Contains(ctx.Request.URL.Path, "auth") {
			ctx.Next()
			return
		}*/

		token_header := ctx.GetHeader("Authorization")
		token := strings.Split(token_header, " ")
		if len(token) <= 1 || token[1] == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			ctx.Abort()
			return
		}

		user_id, is_valid := services.ValidateToken(token[1])
		if is_valid {
			ctx.Set("user_id", user_id)
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
