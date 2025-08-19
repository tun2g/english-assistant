package video

import (
	"context"
	"fmt"
	"strings"

	"app-backend/internal/types"
	"app-backend/pkg/gemini"
	"app-backend/pkg/youtube"
	"go.uber.org/zap"
)

// Service orchestrates video operations across different providers
type Service struct {
	providers   map[types.VideoProvider]ProviderServiceInterface
	translator  *gemini.Service
	logger      *zap.Logger
}

// Config holds configuration for the video service
type Config struct {
	YouTubeAPIKey string
	GeminiAPIKey  string
	Logger        *zap.Logger
}

// NewService creates a new video service with all providers
func NewService(config *Config) (*Service, error) {
	service := &Service{
		providers: make(map[types.VideoProvider]ProviderServiceInterface),
		logger:    config.Logger,
	}

	// Initialize YouTube service
	if config.YouTubeAPIKey != "" {
		youtubeService := youtube.NewService(config.YouTubeAPIKey, config.Logger)
		service.providers[types.ProviderYouTube] = youtubeService
	}

	// Initialize translation service
	if config.GeminiAPIKey != "" {
		translator := gemini.NewService(config.GeminiAPIKey, config.Logger)
		service.translator = translator
	}

	return service, nil
}

// NewVideoService creates a new video service with initialized services (for container injection)
func NewVideoService(youtubeService *youtube.Service, geminiService *gemini.Service, logger *zap.Logger) ServiceInterface {
	service := &Service{
		providers: make(map[types.VideoProvider]ProviderServiceInterface),
		logger:    logger,
	}

	if youtubeService != nil {
		service.providers[types.ProviderYouTube] = youtubeService
	}

	if geminiService != nil {
		service.translator = geminiService
	}

	return service
}

// Close closes all services
func (s *Service) Close() error {
	if s.translator != nil {
		if err := s.translator.Close(); err != nil {
			s.logger.Error("Failed to close translator", zap.Error(err))
		}
	}
	return nil
}

// DetectProvider detects the video provider from URL or video ID
func (s *Service) DetectProvider(videoURL string) (types.VideoProvider, string, error) {
	videoURL = strings.TrimSpace(videoURL)
	
	// Check if it's a YouTube URL or video ID
	if s.isYouTubeURL(videoURL) {
		videoID := s.extractYouTubeVideoID(videoURL)
		if videoID != "" {
			return types.ProviderYouTube, videoID, nil
		}
	}

	// Check if it's a direct YouTube video ID
	if provider, ok := s.providers[types.ProviderYouTube]; ok {
		if provider.ValidateVideoID(videoURL) {
			return types.ProviderYouTube, videoURL, nil
		}
	}

	return "", "", fmt.Errorf("unsupported video provider or invalid URL: %s", videoURL)
}

// GetVideoInfo retrieves video information
func (s *Service) GetVideoInfo(ctx context.Context, provider types.VideoProvider, videoID string) (*types.VideoInfo, error) {
	service, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported", provider)
	}

	return service.GetVideoInfo(ctx, videoID)
}

// GetTranscript retrieves video transcript
func (s *Service) GetTranscript(ctx context.Context, provider types.VideoProvider, videoID string, language string) (*types.Transcript, error) {
	service, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported", provider)
	}

	return service.GetTranscript(ctx, videoID, language)
}

// GetDualLanguageTranscript retrieves transcript and translates it
func (s *Service) GetDualLanguageTranscript(ctx context.Context, provider types.VideoProvider, videoID string, sourceLang string, targetLang string) (*types.DualLanguageTranscript, error) {
	if s.translator == nil {
		return nil, fmt.Errorf("translation service not available")
	}

	// Get original transcript
	transcript, err := s.GetTranscript(ctx, provider, videoID, sourceLang)
	if err != nil {
		return nil, fmt.Errorf("failed to get transcript: %w", err)
	}

	if !transcript.Available || len(transcript.Segments) == 0 {
		return &types.DualLanguageTranscript{
			VideoID:  videoID,
			Provider: provider,
			Segments: []types.TranscriptSegment{},
		}, nil
	}

	// Detect source language if not provided
	detectedSourceLang := transcript.Language
	if sourceLang == "" && len(transcript.Segments) > 0 {
		// Use first few segments to detect language
		sampleText := ""
		for i, segment := range transcript.Segments {
			if i >= 3 { // Use first 3 segments for detection
				break
			}
			sampleText += segment.Text + " "
		}

		if detectedLang, err := s.translator.DetectLanguage(ctx, sampleText); err == nil {
			detectedSourceLang = detectedLang
		}
	}

	// Translate segments
	translations, err := s.translator.TranslateSegments(ctx, transcript.Segments, targetLang, detectedSourceLang)
	if err != nil {
		return nil, fmt.Errorf("failed to translate segments: %w", err)
	}

	return &types.DualLanguageTranscript{
		VideoID:      videoID,
		Provider:     provider,
		SourceLang:   detectedSourceLang,
		TargetLang:   targetLang,
		Segments:     transcript.Segments,
		Translations: translations,
		Cached:       false, // TODO: implement caching
	}, nil
}

// GetAvailableLanguages returns available transcript languages
func (s *Service) GetAvailableLanguages(ctx context.Context, provider types.VideoProvider, videoID string) ([]types.Language, error) {
	service, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported", provider)
	}

	return service.GetAvailableLanguages(ctx, videoID)
}

// GetCapabilities returns video capabilities
func (s *Service) GetCapabilities(ctx context.Context, provider types.VideoProvider, videoID string) (*types.VideoCapabilities, error) {
	service, ok := s.providers[provider]
	if !ok {
		return nil, fmt.Errorf("provider %s not supported", provider)
	}

	return service.GetCapabilities(ctx, videoID)
}

// GetSupportedProviders returns list of supported providers
func (s *Service) GetSupportedProviders() []types.VideoProvider {
	var providers []types.VideoProvider
	for provider := range s.providers {
		providers = append(providers, provider)
	}
	return providers
}

// GetSupportedLanguages returns list of supported translation languages
func (s *Service) GetSupportedLanguages() []types.Language {
	if s.translator == nil {
		return []types.Language{}
	}
	return s.translator.GetSupportedLanguages()
}

// isYouTubeURL checks if the URL is a YouTube URL
func (s *Service) isYouTubeURL(url string) bool {
	return strings.Contains(url, "youtube.com") || 
		   strings.Contains(url, "youtu.be") || 
		   strings.Contains(url, "youtube-nocookie.com")
}

// extractYouTubeVideoID extracts video ID from YouTube URL
func (s *Service) extractYouTubeVideoID(url string) string {
	// Handle different YouTube URL formats
	// https://www.youtube.com/watch?v=VIDEO_ID
	// https://youtu.be/VIDEO_ID
	// https://www.youtube.com/embed/VIDEO_ID
	
	if strings.Contains(url, "watch?v=") {
		parts := strings.Split(url, "watch?v=")
		if len(parts) > 1 {
			videoID := strings.Split(parts[1], "&")[0]
			return videoID
		}
	}
	
	if strings.Contains(url, "youtu.be/") {
		parts := strings.Split(url, "youtu.be/")
		if len(parts) > 1 {
			videoID := strings.Split(parts[1], "?")[0]
			return videoID
		}
	}
	
	if strings.Contains(url, "embed/") {
		parts := strings.Split(url, "embed/")
		if len(parts) > 1 {
			videoID := strings.Split(parts[1], "?")[0]
			return videoID
		}
	}
	
	return ""
}