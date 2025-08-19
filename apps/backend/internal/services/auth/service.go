package auth

import (
	"app-backend/internal/dto"
	"app-backend/internal/errors"
	"app-backend/internal/models"
	"app-backend/internal/repositories"
	"app-backend/internal/services/jwt"
	"app-backend/internal/services/user"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	userService user.ServiceInterface
	jwtService  jwt.ServiceInterface
	sessionRepo repositories.SessionRepositoryInterface
}

func NewAuthService(
	userService user.ServiceInterface,
	jwtService jwt.ServiceInterface,
	sessionRepo repositories.SessionRepositoryInterface,
) ServiceInterface {
	return &Service{
		userService: userService,
		jwtService:  jwtService,
		sessionRepo: sessionRepo,
	}
}

func (s *Service) Register(req *dto.RegisterRequest, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	// Create user
	user, err := s.userService.CreateUser(req)
	if err != nil {
		return nil, err
	}

	// Generate tokens and create session
	return s.createAuthResponse(user, ipAddress, userAgent)
}

func (s *Service) Login(req *dto.LoginRequest, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	// Get user by email
	user, err := s.userService.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.NewAppError("Invalid credentials", nil, http.StatusUnauthorized)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.NewAppError("Account is disabled", nil, http.StatusUnauthorized)
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.NewAppError("Invalid credentials", nil, http.StatusUnauthorized)
	}

	// Generate tokens and create session
	return s.createAuthResponse(user, ipAddress, userAgent)
}

func (s *Service) Logout(userID uint, sessionID uint) error {
	// Deactivate the specific session
	err := s.sessionRepo.DeactivateSession(sessionID)
	if err != nil {
		return errors.NewAppError("Failed to logout", err, http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) LogoutAll(userID uint) error {
	// Deactivate all user sessions
	err := s.sessionRepo.DeactivateUserSessions(userID)
	if err != nil {
		return errors.NewAppError("Failed to logout from all devices", err, http.StatusInternalServerError)
	}
	return nil
}

func (s *Service) RefreshToken(req *dto.RefreshTokenRequest, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewAppError("Invalid refresh token", err, http.StatusUnauthorized)
	}

	// Check if it's a refresh token
	if claims.TokenType != "refresh" {
		return nil, errors.NewAppError("Invalid token type", nil, http.StatusUnauthorized)
	}

	// Get session by token hash
	tokenHash := s.jwtService.GetTokenHash(req.RefreshToken)
	session, err := s.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("Session not found", nil, http.StatusUnauthorized)
		}
		return nil, errors.NewAppError("Failed to validate session", err, http.StatusInternalServerError)
	}

	// Check if session is active and not expired
	if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewAppError("Session expired", nil, http.StatusUnauthorized)
	}

	// Get user
	user, err := s.userService.GetUser(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, errors.NewAppError("Account is disabled", nil, http.StatusUnauthorized)
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, errors.NewAppError("Failed to generate access token", err, http.StatusInternalServerError)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, errors.NewAppError("Failed to generate refresh token", err, http.StatusInternalServerError)
	}

	// Update session
	session.TokenHash = s.jwtService.GetTokenHash(refreshToken)
	session.LastUsed = time.Now()
	session.ExpiresAt = time.Now().Add(s.jwtService.GetRefreshTokenTTL())
	session.IPAddress = ipAddress
	session.UserAgent = userAgent

	err = s.sessionRepo.Update(session)
	if err != nil {
		return nil, errors.NewAppError("Failed to update session", err, http.StatusInternalServerError)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtService.GetAccessTokenTTL().Seconds()),
		User: &dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *Service) ValidateSession(tokenHash string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("Session not found", nil, http.StatusUnauthorized)
		}
		return nil, errors.NewAppError("Failed to validate session", err, http.StatusInternalServerError)
	}

	// Check if session is active and not expired
	if !session.IsActive || session.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewAppError("Session expired", nil, http.StatusUnauthorized)
	}

	// Update last used timestamp
	err = s.sessionRepo.UpdateLastUsed(session.ID)
	if err != nil {
		// Log error but don't fail the request
		// logger.Error("Failed to update session last used", "error", err)
	}

	return session, nil
}

func (s *Service) GetUserSessions(userID uint) ([]*dto.SessionResponse, error) {
	sessions, err := s.sessionRepo.GetActiveSessionsByUserID(userID)
	if err != nil {
		return nil, errors.NewAppError("Failed to get user sessions", err, http.StatusInternalServerError)
	}

	sessionResponses := make([]*dto.SessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = &dto.SessionResponse{
			ID:        session.ID,
			UserAgent: session.UserAgent,
			IPAddress: session.IPAddress,
			LastUsed:  session.LastUsed,
			ExpiresAt: session.ExpiresAt,
			IsActive:  session.IsActive,
			CreatedAt: session.CreatedAt,
		}
	}

	return sessionResponses, nil
}

func (s *Service) RevokeSession(userID uint, sessionID uint) error {
	// Verify the session belongs to the user
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError("Session not found", nil, http.StatusNotFound)
		}
		return errors.NewAppError("Failed to get session", err, http.StatusInternalServerError)
	}

	if session.UserID != userID {
		return errors.NewAppError("Session does not belong to user", nil, http.StatusForbidden)
	}

	// Deactivate the session
	err = s.sessionRepo.DeactivateSession(sessionID)
	if err != nil {
		return errors.NewAppError("Failed to revoke session", err, http.StatusInternalServerError)
	}

	return nil
}

func (s *Service) createAuthResponse(user *models.User, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	// Create session first (without token hash)
	session := &models.Session{
		UserID:    user.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		IsActive:  true,
		LastUsed:  time.Now(),
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTokenTTL()),
	}

	err := s.sessionRepo.Create(session)
	if err != nil {
		return nil, errors.NewAppError("Failed to create session", err, http.StatusInternalServerError)
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, errors.NewAppError("Failed to generate access token", err, http.StatusInternalServerError)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Email, user.Role, session.ID)
	if err != nil {
		return nil, errors.NewAppError("Failed to generate refresh token", err, http.StatusInternalServerError)
	}

	// Update session with token hash
	session.TokenHash = s.jwtService.GetTokenHash(refreshToken)
	err = s.sessionRepo.Update(session)
	if err != nil {
		return nil, errors.NewAppError("Failed to update session with token hash", err, http.StatusInternalServerError)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtService.GetAccessTokenTTL().Seconds()),
		User: &dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}