package routes

import (
	"app-backend/internal/handlers/translation"

	"github.com/gin-gonic/gin"
)

// SetupTranslationRoutes configures translation-related routes
func SetupTranslationRoutes(rg *gin.RouterGroup, handler translation.HandlerInterface) {
	translationGroup := rg.Group("/translate")
	{
		// Text translation endpoint
		translationGroup.POST("", handler.TranslateTexts)
	}
}