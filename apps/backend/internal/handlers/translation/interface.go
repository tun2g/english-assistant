package translation

import "github.com/gin-gonic/gin"

// HandlerInterface defines the contract for translation HTTP handlers
type HandlerInterface interface {
	// TranslateTexts handles text translation requests
	TranslateTexts(c *gin.Context)
}