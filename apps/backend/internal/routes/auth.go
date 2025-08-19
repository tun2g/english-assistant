package routes

import (
	"app-backend/internal/handlers/auth"
	"app-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures all authentication routes
func SetupAuthRoutes(router *gin.RouterGroup, authHandler auth.HandlerInterface, authMiddleware *middleware.AuthMiddleware) {
	authGroup := router.Group("/auth")
	{
		// Public routes (no authentication required)
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)

		// Protected routes (authentication required)
		protected := authGroup.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/logout-all", authHandler.LogoutAll)
			protected.GET("/sessions", authHandler.GetSessions)
			protected.DELETE("/sessions/:sessionId", authHandler.RevokeSession)
		}
	}
}