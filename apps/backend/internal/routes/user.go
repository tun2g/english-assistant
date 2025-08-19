package routes

import (
	"app-backend/internal/handlers/user"
	"app-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user management routes
func SetupUserRoutes(router *gin.RouterGroup, userHandler user.HandlerInterface, authMiddleware *middleware.AuthMiddleware) {
	userGroup := router.Group("/user")
	userGroup.Use(authMiddleware.RequireAuth()) // All user routes require authentication
	{
		// User profile management
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.POST("/change-password", userHandler.ChangePassword)
		userGroup.DELETE("/account", userHandler.DeleteAccount)

		// Admin only routes
		adminGroup := userGroup.Group("")
		adminGroup.Use(authMiddleware.RequireRole("admin"))
		{
			adminGroup.GET("/list", userHandler.ListUsers)
		}
	}
}