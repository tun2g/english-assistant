package routes

import (
	"app-backend/internal/handlers/oauth"
	
	"github.com/gin-gonic/gin"
)

// SetupOAuthRoutes sets up all OAuth related routes
func SetupOAuthRoutes(rg *gin.RouterGroup, handler oauth.HandlerInterface) {
	oauthGroup := rg.Group("/oauth")
	{
		// YouTube OAuth routes
		youtube := oauthGroup.Group("/youtube")
		{
			// Initiate YouTube OAuth flow
			youtube.GET("/auth", handler.InitiateYouTubeAuth)
			
			// Handle YouTube OAuth callback
			youtube.GET("/callback", handler.HandleYouTubeCallback)
			
			// Get current authentication status
			youtube.GET("/status", handler.GetAuthStatus)
			
			// Revoke current authentication
			youtube.POST("/revoke", handler.RevokeYouTubeAuth)
		}
	}
}