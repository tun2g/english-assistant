package docs

import (
	"app-backend/internal/dto"
)

// NewVideoDocs creates instances of video-related DTOs for swagger documentation
// This function is never called but ensures the DTOs are considered "used" by the linter
func NewVideoDocs() {
	_ = dto.VideoInfoRequest{}
	_ = dto.VideoInfoResponse{}
	_ = dto.GetTranscriptRequest{}
	_ = dto.GetTranscriptResponse{}
	_ = dto.GetAvailableLanguagesResponse{}
	_ = dto.VideoCapabilitiesResponse{}
	_ = dto.GetSupportedProvidersResponse{}
	_ = dto.GetSupportedLanguagesResponse{}
}

// VideoGetInfo godoc
// @Summary Get video information
// @Description Get basic information about a video from any supported provider
// @Tags video
// @Accept json
// @Produce json
// @Param videoUrl path string true "Video URL (base64 encoded)"
// @Success 200 {object} dto.VideoInfoResponse "Video information"
// @Failure 400 {object} dto.ErrorResponse "Invalid video URL"
// @Failure 404 {object} dto.ErrorResponse "Video not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/video/{videoUrl}/info [get]
// @Security BearerAuth
func VideoGetInfo() {}

// VideoGetTranscript godoc
// @Summary Get video transcript
// @Description Get transcript for a video in the specified language
// @Tags video
// @Accept json
// @Produce json
// @Param videoUrl path string true "Video URL (base64 encoded)"
// @Param language query string false "Language code (e.g., 'en', 'es')" default(en)
// @Success 200 {object} dto.GetTranscriptResponse "Video transcript"
// @Failure 400 {object} dto.ErrorResponse "Invalid parameters"
// @Failure 404 {object} dto.ErrorResponse "Transcript not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/video/{videoUrl}/transcript [get]
// @Security BearerAuth
func VideoGetTranscript() {}

// VideoGetAvailableLanguages godoc
// @Summary Get available transcript languages
// @Description Get list of available transcript languages for a video
// @Tags video
// @Accept json
// @Produce json
// @Param videoUrl path string true "Video URL (base64 encoded)"
// @Success 200 {object} dto.GetAvailableLanguagesResponse "Available languages"
// @Failure 400 {object} dto.ErrorResponse "Invalid video URL"
// @Failure 404 {object} dto.ErrorResponse "Video not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/video/{videoUrl}/languages [get]
// @Security BearerAuth
func VideoGetAvailableLanguages() {}

// VideoGetCapabilities godoc
// @Summary Get video capabilities
// @Description Get capabilities and features available for a video
// @Tags video
// @Accept json
// @Produce json
// @Param videoUrl path string true "Video URL (base64 encoded)"
// @Success 200 {object} dto.VideoCapabilitiesResponse "Video capabilities"
// @Failure 400 {object} dto.ErrorResponse "Invalid video URL"
// @Failure 404 {object} dto.ErrorResponse "Video not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/v1/video/{videoUrl}/capabilities [get]
// @Security BearerAuth
func VideoGetCapabilities() {}

// VideoGetSupportedProviders godoc
// @Summary Get supported video providers
// @Description Get list of supported video providers and their capabilities
// @Tags video
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetSupportedProvidersResponse "Supported providers"
// @Router /api/v1/video/providers [get]
// @Security BearerAuth
func VideoGetSupportedProviders() {}

// VideoGetSupportedLanguages godoc
// @Summary Get supported translation languages
// @Description Get list of supported languages for AI translation
// @Tags video
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetSupportedLanguagesResponse "Supported translation languages"
// @Router /api/v1/video/languages [get]
// @Security BearerAuth
func VideoGetSupportedLanguages() {}