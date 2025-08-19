package auth

import (
	"app-backend/internal/dto"
	"app-backend/internal/errors"
	"app-backend/internal/logger"
	"app-backend/internal/services/auth"
	"app-backend/internal/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	authService auth.ServiceInterface
	logger      *logger.Logger
}

func NewAuthHandler(authService auth.ServiceInterface, logger *logger.Logger) HandlerInterface {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authService.Register(&req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Registration failed", zap.Error(err), zap.String("email", req.Email))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected registration error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("User registered successfully", zap.Uint("user_id", response.User.ID), zap.String("email", response.User.Email))
	c.JSON(http.StatusCreated, response)
}

func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authService.Login(&req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Login failed", zap.Error(err), zap.String("email", req.Email))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected login error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("User logged in successfully", zap.Uint("user_id", response.User.ID), zap.String("email", response.User.Email))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) Logout(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.authService.Logout(userCtx.UserID, userCtx.SessionID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Logout failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected logout error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("User logged out successfully", zap.Uint("user_id", userCtx.UserID), zap.Uint("session_id", userCtx.SessionID))
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) LogoutAll(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.authService.LogoutAll(userCtx.UserID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Logout all failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected logout all error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("User logged out from all sessions", zap.Uint("user_id", userCtx.UserID))
	c.JSON(http.StatusOK, gin.H{"message": "Logged out from all sessions successfully"})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid refresh token request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	response, err := h.authService.RefreshToken(&req, ipAddress, userAgent)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Token refresh failed", zap.Error(err))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected token refresh error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Token refreshed successfully", zap.Uint("user_id", response.User.ID))
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetSessions(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessions, err := h.authService.GetUserSessions(userCtx.UserID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Get sessions failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected get sessions error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *Handler) RevokeSession(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	sessionIDStr := c.Param("sessionId")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	err = h.authService.RevokeSession(userCtx.UserID, uint(sessionID))
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Revoke session failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID), zap.Uint64("session_id", sessionID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected revoke session error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Session revoked successfully", zap.Uint("user_id", userCtx.UserID), zap.Uint64("session_id", sessionID))
	c.JSON(http.StatusOK, gin.H{"message": "Session revoked successfully"})
}