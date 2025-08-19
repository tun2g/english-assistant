package oauth

import "github.com/gin-gonic/gin"

// HandlerInterface defines the interface for OAuth HTTP handlers
type HandlerInterface interface {
	// InitiateYouTubeAuth starts the YouTube OAuth flow
	InitiateYouTubeAuth(c *gin.Context)
	
	// HandleYouTubeCallback handles the OAuth callback from YouTube
	HandleYouTubeCallback(c *gin.Context)
	
	// GetAuthStatus checks the current YouTube authentication status
	GetAuthStatus(c *gin.Context)
	
	// RevokeYouTubeAuth revokes the current YouTube authentication
	RevokeYouTubeAuth(c *gin.Context)
}