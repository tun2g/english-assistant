package auth

import (
	"app-backend/internal/dto"
	"app-backend/internal/models"
)

type ServiceInterface interface {
	Register(req *dto.RegisterRequest, ipAddress, userAgent string) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest, ipAddress, userAgent string) (*dto.AuthResponse, error)
	Logout(userID uint, sessionID uint) error
	LogoutAll(userID uint) error
	RefreshToken(req *dto.RefreshTokenRequest, ipAddress, userAgent string) (*dto.AuthResponse, error)
	ValidateSession(tokenHash string) (*models.Session, error)
	GetUserSessions(userID uint) ([]*dto.SessionResponse, error)
	RevokeSession(userID uint, sessionID uint) error
}