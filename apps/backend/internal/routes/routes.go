package routes

import (
	"app-backend/internal/handlers/auth"
	"app-backend/internal/handlers/oauth"
	"app-backend/internal/handlers/translation"
	"app-backend/internal/handlers/user"
	"app-backend/internal/handlers/video"
	"app-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RouteConfig holds all the dependencies needed for route setup
type RouteConfig struct {
	AuthHandler        auth.HandlerInterface
	UserHandler        user.HandlerInterface
	VideoHandler       video.HandlerInterface
	OAuthHandler       oauth.HandlerInterface
	TranslationHandler translation.HandlerInterface
	AuthMiddleware     *middleware.AuthMiddleware
}

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, config *RouteConfig) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API version 1 routes
	v1 := router.Group("/api/v1")
	{
		// Setup all route groups
		SetupAuthRoutes(v1, config.AuthHandler, config.AuthMiddleware)
		SetupUserRoutes(v1, config.UserHandler, config.AuthMiddleware)
		SetupVideoRoutes(v1, config.VideoHandler, config.AuthMiddleware)
		SetupTranslationRoutes(v1, config.TranslationHandler)
		SetupOAuthRoutes(v1, config.OAuthHandler)
	}

	// Setup Swagger documentation routes
	SetupSwaggerRoutes(router)
}