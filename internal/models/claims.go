package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	Role     int    `json:"role"`
	Status   int    `json:"status"`
	jwt.RegisteredClaims
}
