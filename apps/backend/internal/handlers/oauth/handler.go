package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"app-backend/internal/dto"
	"app-backend/internal/logger"
	oauthService "app-backend/internal/services/oauth"
	
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler implements OAuth HTTP handlers
type Handler struct {
	youtubeOAuth oauthService.ServiceInterface
	logger       *logger.Logger
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(youtubeOAuth oauthService.ServiceInterface, logger *logger.Logger) HandlerInterface {
	return &Handler{
		youtubeOAuth: youtubeOAuth,
		logger:       logger,
	}
}

// InitiateYouTubeAuth starts the YouTube OAuth flow
func (h *Handler) InitiateYouTubeAuth(c *gin.Context) {
	// Generate random state for security
	state := h.generateRandomState()
	
	// Store state in memory/session for verification (instead of cookie)
	// For Chrome extension OAuth, cookies are not reliable due to cross-origin restrictions
	h.youtubeOAuth.StoreState(state)
	
	// Generate authorization URL
	authURL := h.youtubeOAuth.GenerateAuthURL(state)
	
	h.logger.Info("Initiating YouTube OAuth flow", zap.String("state", state))
	
	c.JSON(http.StatusOK, gin.H{
		"authUrl": authURL,
		"state":   state,
	})
}

// HandleYouTubeCallback handles the OAuth callback from YouTube
func (h *Handler) HandleYouTubeCallback(c *gin.Context) {
	// Get authorization code and state from query parameters
	code := c.Query("code")
	state := c.Query("state")
	
	if code == "" {
		h.logger.Error("No authorization code received in callback")
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Authorization code not provided",
		})
		return
	}
	
	// Verify state parameter to prevent CSRF attacks
	if !h.youtubeOAuth.ValidateAndClearState(state) {
		h.logger.Error("Invalid OAuth state", zap.String("received", state))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid state parameter",
		})
		return
	}
	
	// Exchange code for tokens
	token, err := h.youtubeOAuth.ExchangeCodeForTokens(c.Request.Context(), code)
	if err != nil {
		h.logger.Error("Failed to exchange code for tokens", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to complete OAuth flow",
			Details: err.Error(),
		})
		return
	}
	
	h.logger.Info("Successfully completed YouTube OAuth flow")
	
	// For web flow, redirect to success page or return success response
	if c.Query("redirect") == "web" {
		// Redirect to frontend success page
		c.Redirect(http.StatusFound, "/oauth/success")
		return
	}
	
	// For API flow, return token info (without sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "YouTube authentication completed",
		"expiresAt": token.Expiry,
	})
}

// GetAuthStatus checks the current YouTube authentication status
func (h *Handler) GetAuthStatus(c *gin.Context) {
	isAuthenticated := h.youtubeOAuth.IsAuthenticated()
	
	response := gin.H{
		"authenticated": isAuthenticated,
	}
	
	// If authenticated, get token expiry info
	if isAuthenticated {
		token, err := h.youtubeOAuth.LoadToken()
		if err == nil {
			response["expiresAt"] = token.Expiry
			response["valid"] = token.Valid()
		}
	}
	
	c.JSON(http.StatusOK, response)
}

// RevokeYouTubeAuth revokes the current YouTube authentication
func (h *Handler) RevokeYouTubeAuth(c *gin.Context) {
	if !h.youtubeOAuth.IsAuthenticated() {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "No active authentication to revoke",
		})
		return
	}
	
	err := h.youtubeOAuth.RevokeToken(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to revoke YouTube authentication", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to revoke authentication",
			Details: err.Error(),
		})
		return
	}
	
	h.logger.Info("Successfully revoked YouTube authentication")
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "YouTube authentication revoked",
	})
}

// generateRandomState generates a random state string for OAuth flow
func (h *Handler) generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}