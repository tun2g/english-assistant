package video

import (
	"net/http"
	"net/url"
	"sync"

	"app-backend/internal/dto"
	"app-backend/internal/logger"
	"app-backend/internal/services/transcript"
	"app-backend/internal/services/transcript/types"
	"app-backend/internal/services/video"
	internalTypes "app-backend/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler implements video HTTP handlers
type Handler struct {
	videoService      video.ServiceInterface
	transcriptService transcript.ServiceInterface
	logger            *logger.Logger
}

// NewVideoHandler creates a new video handler
func NewVideoHandler(videoService video.ServiceInterface, transcriptService transcript.ServiceInterface, logger *logger.Logger) HandlerInterface {
	return &Handler{
		videoService:      videoService,
		transcriptService: transcriptService,
		logger:            logger,
	}
}

// GetVideoInfo retrieves basic information about a video
func (h *Handler) GetVideoInfo(c *gin.Context) {
	var req dto.VideoInfoRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL",
			Details: err.Error(),
		})
		return
	}

	// URL decode the video URL
	decodedURL, err := url.QueryUnescape(req.VideoURL)
	if err != nil {
		h.logger.Error("Failed to decode URL", zap.String("url", req.VideoURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL format",
			Details: err.Error(),
		})
		return
	}

	// Detect provider and extract video ID
	provider, videoID, err := h.videoService.DetectProvider(decodedURL)
	if err != nil {
		h.logger.Error("Failed to detect provider", zap.String("url", decodedURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Unsupported video provider or invalid URL",
			Details: err.Error(),
		})
		return
	}

	// Fetch video info and capabilities concurrently for better performance
	var videoInfo *internalTypes.VideoInfo
	var capabilities *internalTypes.VideoCapabilities
	var videoErr, capErr error
	var wg sync.WaitGroup

	wg.Add(2)

	// Fetch video info in parallel
	go func() {
		defer wg.Done()
		videoInfo, videoErr = h.videoService.GetVideoInfo(c.Request.Context(), provider, videoID)
	}()

	// Fetch capabilities in parallel
	go func() {
		defer wg.Done()
		capabilities, capErr = h.videoService.GetCapabilities(c.Request.Context(), provider, videoID)
	}()

	wg.Wait()

	// Check for critical video info error
	if videoErr != nil {
		h.logger.Error("Failed to get video info", 
			zap.String("provider", string(provider)),
			zap.String("videoID", videoID),
			zap.Error(videoErr))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to retrieve video information",
			Details: videoErr.Error(),
		})
		return
	}

	// Capabilities error is non-fatal (existing behavior)
	if capErr != nil {
		h.logger.Warn("Failed to get capabilities", zap.Error(capErr))
		capabilities = nil
	}

	response := dto.VideoInfoResponse{
		ID:           videoInfo.ID,
		Provider:     videoInfo.Provider,
		Title:        videoInfo.Title,
		Description:  videoInfo.Description,
		Duration:     videoInfo.Duration,
		ThumbnailURL: videoInfo.ThumbnailURL,
		URL:          videoInfo.URL,
	}

	if capabilities != nil {
		capResponse := dto.ConvertToVideoCapabilitiesResponse(*capabilities)
		response.Capabilities = &capResponse
	}

	c.JSON(http.StatusOK, response)
}

// GetTranscript retrieves transcript for a video
func (h *Handler) GetTranscript(c *gin.Context) {
	var req dto.GetTranscriptRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL",
			Details: err.Error(),
		})
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("Invalid query parameters", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// URL decode the video URL
	decodedURL, err := url.QueryUnescape(req.VideoURL)
	if err != nil {
		h.logger.Error("Failed to decode URL", zap.String("url", req.VideoURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL format",
			Details: err.Error(),
		})
		return
	}

	// Create transcript request
	transcriptReq := &types.TranscriptRequest{
		VideoURL: decodedURL,
		Language: req.Language,
	}

	// Get transcript using our new transcript service
	transcript, err := h.transcriptService.GetTranscript(c.Request.Context(), transcriptReq)
	if err != nil {
		h.logger.Error("Failed to get transcript",
			zap.String("video_url", decodedURL),
			zap.String("language", req.Language),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to retrieve transcript",
			Details: err.Error(),
		})
		return
	}

	// Convert to response format
	var segments []dto.TranscriptSegmentResponse
	for i, segment := range transcript.Segments {
		segmentResponse := dto.ConvertFromTranscriptServiceSegment(segment)
		segmentResponse.Index = i + 1 // Set proper index
		segments = append(segments, segmentResponse)
	}

	response := dto.GetTranscriptResponse{
		VideoID:   transcript.VideoID,
		Provider:  internalTypes.VideoProvider(transcript.Provider),
		Language:  transcript.Language,
		Segments:  segments,
		Available: true, // If we got here, transcript is available
		Source:    transcript.Provider,
	}

	c.JSON(http.StatusOK, response)
}


// GetAvailableLanguages returns available transcript languages for a video
func (h *Handler) GetAvailableLanguages(c *gin.Context) {
	var req dto.GetAvailableLanguagesRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL",
			Details: err.Error(),
		})
		return
	}

	// URL decode the video URL
	decodedURL, err := url.QueryUnescape(req.VideoURL)
	if err != nil {
		h.logger.Error("Failed to decode URL", zap.String("url", req.VideoURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL format",
			Details: err.Error(),
		})
		return
	}

	// Detect provider and extract video ID
	provider, videoID, err := h.videoService.DetectProvider(decodedURL)
	if err != nil {
		h.logger.Error("Failed to detect provider", zap.String("url", decodedURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Unsupported video provider or invalid URL",
			Details: err.Error(),
		})
		return
	}

	// Get available languages
	languages, err := h.videoService.GetAvailableLanguages(c.Request.Context(), provider, videoID)
	if err != nil {
		h.logger.Error("Failed to get available languages",
			zap.String("provider", string(provider)),
			zap.String("videoID", videoID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to retrieve available languages",
			Details: err.Error(),
		})
		return
	}

	// Convert to response format
	var languageResponses []dto.LanguageResponse
	for _, lang := range languages {
		languageResponses = append(languageResponses, dto.ConvertToLanguageResponse(lang))
	}

	response := dto.GetAvailableLanguagesResponse{
		VideoID:   videoID,
		Provider:  provider,
		Languages: languageResponses,
	}

	c.JSON(http.StatusOK, response)
}

// GetCapabilities returns capabilities for a video
func (h *Handler) GetCapabilities(c *gin.Context) {
	var req dto.GetAvailableLanguagesRequest
	if err := c.ShouldBindUri(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL",
			Details: err.Error(),
		})
		return
	}

	// URL decode the video URL
	decodedURL, err := url.QueryUnescape(req.VideoURL)
	if err != nil {
		h.logger.Error("Failed to decode URL", zap.String("url", req.VideoURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Invalid video URL format",
			Details: err.Error(),
		})
		return
	}

	// Detect provider and extract video ID
	provider, videoID, err := h.videoService.DetectProvider(decodedURL)
	if err != nil {
		h.logger.Error("Failed to detect provider", zap.String("url", decodedURL), zap.Error(err))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Unsupported video provider or invalid URL",
			Details: err.Error(),
		})
		return
	}

	// Get capabilities
	capabilities, err := h.videoService.GetCapabilities(c.Request.Context(), provider, videoID)
	if err != nil {
		h.logger.Error("Failed to get capabilities",
			zap.String("provider", string(provider)),
			zap.String("videoID", videoID),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to retrieve video capabilities",
			Details: err.Error(),
		})
		return
	}

	response := dto.ConvertToVideoCapabilitiesResponse(*capabilities)
	c.JSON(http.StatusOK, response)
}

// GetSupportedProviders returns list of supported video providers
func (h *Handler) GetSupportedProviders(c *gin.Context) {
	providers := h.videoService.GetSupportedProviders()
	
	response := dto.GetSupportedProvidersResponse{
		Providers: providers,
	}

	c.JSON(http.StatusOK, response)
}

// GetSupportedLanguages returns list of supported translation languages
func (h *Handler) GetSupportedLanguages(c *gin.Context) {
	languages := h.videoService.GetSupportedLanguages()
	
	var languageResponses []dto.LanguageResponse
	for _, lang := range languages {
		languageResponses = append(languageResponses, dto.ConvertToLanguageResponse(lang))
	}

	response := dto.GetSupportedLanguagesResponse{
		Languages: languageResponses,
	}

	c.JSON(http.StatusOK, response)
}