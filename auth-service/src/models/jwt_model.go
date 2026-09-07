package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	UserId int       `json:"sub"`
	Name   string    `json:"name"`
	Roles  *[]string `json:"roles"`
	Exp    int64     `json:"exp"`
	jwt.RegisteredClaims
}
