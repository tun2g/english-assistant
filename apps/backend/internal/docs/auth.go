package docs

import (
	"app-backend/internal/dto"
	"github.com/gin-gonic/gin"
)

// NewAuthDocs creates instances of auth-related DTOs for swagger documentation
// This function is never called but ensures the DTOs are considered "used" by the linter
func NewAuthDocs() {
	_ = dto.RegisterRequest{}
	_ = dto.LoginRequest{}
	_ = dto.RefreshTokenRequest{}
	_ = dto.ChangePasswordRequest{}
	_ = dto.AuthResponse{}
	_ = dto.SessionResponse{}
}

// AuthRegister godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/register [post]
func AuthRegister(c *gin.Context) {}

// AuthLogin godoc
// @Summary Login user
// @Description Authenticate user and return access and refresh tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func AuthLogin(c *gin.Context) {}

// AuthLogout godoc
// @Summary Logout user
// @Description Logout user from current session
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Successfully logged out"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/logout [post]
func AuthLogout(c *gin.Context) {}

// AuthLogoutAll godoc
// @Summary Logout user from all sessions
// @Description Logout user from all active sessions across all devices
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Successfully logged out from all sessions"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/logout-all [post]
func AuthLogoutAll(c *gin.Context) {}

// AuthRefreshToken godoc
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "Invalid or expired refresh token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/refresh [post]
func AuthRefreshToken(c *gin.Context) {}

// AuthGetSessions godoc
// @Summary Get user sessions
// @Description Get all active sessions for the authenticated user
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} dto.SessionResponse "List of active sessions"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/sessions [get]
func AuthGetSessions(c *gin.Context) {}

// AuthRevokeSession godoc
// @Summary Revoke a specific session
// @Description Revoke a specific session by session ID
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param sessionId path int true "Session ID to revoke"
// @Success 200 {object} map[string]interface{} "Session revoked successfully"
// @Failure 400 {object} map[string]interface{} "Invalid session ID"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Session does not belong to user"
// @Failure 404 {object} map[string]interface{} "Session not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/sessions/{sessionId} [delete]
func AuthRevokeSession(c *gin.Context) {}