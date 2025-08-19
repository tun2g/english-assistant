package jwt

import (
	"time"
)

type ServiceInterface interface {
	GenerateAccessToken(userID uint, email, role string, sessionID uint) (string, error)
	GenerateRefreshToken(userID uint, email, role string, sessionID uint) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GetTokenHash(token string) string
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
}