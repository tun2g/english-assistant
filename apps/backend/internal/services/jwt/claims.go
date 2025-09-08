package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID uint   `json:"session_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}