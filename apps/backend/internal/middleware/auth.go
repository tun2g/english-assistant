package middleware

import (
	"app-backend/internal/errors"
	"app-backend/internal/logger"
	"app-backend/internal/services/auth"
	"app-backend/internal/services/jwt"
	"app-backend/internal/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwtService  jwt.ServiceInterface
	authService auth.ServiceInterface
	logger      *logger.Logger
}

func NewAuthMiddleware(jwtService jwt.ServiceInterface, authService auth.ServiceInterface, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth middleware validates JWT token and sets user context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.logger.Warn("Invalid authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.logger.Warn("Invalid JWT token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if it's an access token
		if claims.TokenType != "access" {
			m.logger.Warn("Invalid token type", zap.String("token_type", claims.TokenType))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Validate session using refresh token hash
		// Note: For access tokens, we don't validate against session directly
		// but we could add additional session validation here if needed

		// Set user context
		userCtx := &types.UserContext{
			UserID:    claims.UserID,
			Email:     claims.Email,
			Role:      claims.Role,
			SessionID: claims.SessionID,
		}
		types.SetUserContext(c, userCtx)

		m.logger.Debug("User authenticated", zap.Uint("user_id", claims.UserID), zap.String("email", claims.Email))
		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, err := types.GetUserContext(c)
		if err != nil {
			m.logger.Error("User context not found", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		if !userCtx.HasRole(roles...) {
			m.logger.Warn("Insufficient permissions", zap.String("user_role", userCtx.Role), zap.Strings("required_roles", roles))
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware validates JWT token if present but doesn't require it
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			c.Next()
			return
		}

		if claims.TokenType == "access" {
			userCtx := &types.UserContext{
				UserID:    claims.UserID,
				Email:     claims.Email,
				Role:      claims.Role,
				SessionID: claims.SessionID,
			}
			types.SetUserContext(c, userCtx)
		}

		c.Next()
	}
}

// ValidateRefreshToken middleware specifically for refresh token endpoints
func (m *AuthMiddleware) ValidateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenReq struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := c.ShouldBindJSON(&tokenReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			c.Abort()
			return
		}

		// Validate refresh token
		claims, err := m.jwtService.ValidateToken(tokenReq.RefreshToken)
		if err != nil {
			m.logger.Warn("Invalid refresh token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
			c.Abort()
			return
		}

		if claims.TokenType != "refresh" {
			m.logger.Warn("Invalid token type for refresh", zap.String("token_type", claims.TokenType))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Validate session
		tokenHash := m.jwtService.GetTokenHash(tokenReq.RefreshToken)
		session, err := m.authService.ValidateSession(tokenHash)
		if err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				c.JSON(appErr.Status, gin.H{"error": appErr.Message})
				c.Abort()
				return
			}
			m.logger.Error("Session validation failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Set context for refresh token processing
		c.Set("user_id", claims.UserID)
		c.Set("session_id", session.ID)
		c.Set("refresh_token", tokenReq.RefreshToken)

		c.Next()
	}
}