package video

import "github.com/gin-gonic/gin"

// HandlerInterface defines the contract for video handlers
type HandlerInterface interface {
	// GetVideoInfo retrieves basic information about a video
	GetVideoInfo(c *gin.Context)
	
	// GetTranscript retrieves transcript for a video
	GetTranscript(c *gin.Context)
	
	// TranslateTranscript translates a video transcript
	TranslateTranscript(c *gin.Context)
	
	// GetAvailableLanguages returns available transcript languages for a video
	GetAvailableLanguages(c *gin.Context)
	
	// GetCapabilities returns capabilities for a video
	GetCapabilities(c *gin.Context)
	
	// GetSupportedProviders returns list of supported video providers
	GetSupportedProviders(c *gin.Context)
	
	// GetSupportedLanguages returns list of supported translation languages
	GetSupportedLanguages(c *gin.Context)
}