package services

import (
	"fmt"
	"log/slog"
	"time"

	"src/models"
	"src/utils"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtSecret = []byte(os.Getenv("JWT_SECRET_TOKEN"))
var jwtSecret = []byte("123456789012345678")

func CreateToken(claims models.JWTClaims, expiration_time int) (string, error) {
	extra_time := utils.GetDefaultIfNull(expiration_time, 60)
	claims.Exp = time.Now().Add(time.Duration(extra_time) * time.Hour).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(token string) (userId int, isValid bool) {
	if claims, ok := validate(token); ok {
		isValid = ok
		userId = claims.UserId
	}
	return
}

func validate(token string) (*models.JWTClaims, bool) {
	// Initialize an empty instance
	claims := &models.JWTClaims{}

	token_obj, err := jwt.ParseWithClaims(token, claims, func(in_token *jwt.Token) (interface{}, error) {
		if _, ok := in_token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", in_token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token_obj.Valid {
		slog.Error("Invalid token. ", "error", err.Error())
		return nil, false
	}

	return claims, true
}
