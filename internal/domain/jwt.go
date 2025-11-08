package domain

import "github.com/golang-jwt/jwt/v5"

// JWTClaims represents the claims stored in JWT token
type JWTClaims struct {
	UserID int32  `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}
