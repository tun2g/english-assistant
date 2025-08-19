package routes

import (
	"app-backend/internal/handlers/video"
	"app-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupVideoRoutes configures video-related routes
func SetupVideoRoutes(rg *gin.RouterGroup, handler video.HandlerInterface, authMiddleware *middleware.AuthMiddleware) {
	videoGroup := rg.Group("/video")
	{
		// Video information and capabilities
		videoGroup.GET("/:videoUrl/info", handler.GetVideoInfo)
		videoGroup.GET("/:videoUrl/capabilities", handler.GetCapabilities)
		
		// Transcript operations
		videoGroup.GET("/:videoUrl/transcript", handler.GetTranscript)
		videoGroup.POST("/:videoUrl/translate", handler.TranslateTranscript)
		videoGroup.GET("/:videoUrl/languages", handler.GetAvailableLanguages)
		
		// System endpoints
		videoGroup.GET("/providers", handler.GetSupportedProviders)
		videoGroup.GET("/languages", handler.GetSupportedLanguages)
	}
}