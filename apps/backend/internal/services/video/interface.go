package video

import (
	"context"
	"app-backend/internal/types"
)

// ServiceInterface defines the contract for the main video service facade
type ServiceInterface interface {
	// DetectProvider detects the video provider from URL or video ID
	DetectProvider(videoURL string) (types.VideoProvider, string, error)
	
	// GetVideoInfo retrieves basic information about a video
	GetVideoInfo(ctx context.Context, provider types.VideoProvider, videoID string) (*types.VideoInfo, error)
	
	// GetTranscript retrieves transcript for a video in specified language
	GetTranscript(ctx context.Context, provider types.VideoProvider, videoID string, language string) (*types.Transcript, error)
	
	// GetAvailableLanguages returns list of available transcript languages
	GetAvailableLanguages(ctx context.Context, provider types.VideoProvider, videoID string) ([]types.Language, error)
	
	// GetCapabilities returns what features are supported for this video
	GetCapabilities(ctx context.Context, provider types.VideoProvider, videoID string) (*types.VideoCapabilities, error)
	
	// GetDualLanguageTranscript retrieves transcript and translates it
	GetDualLanguageTranscript(ctx context.Context, provider types.VideoProvider, videoID string, sourceLang string, targetLang string) (*types.DualLanguageTranscript, error)
	
	// GetSupportedProviders returns list of supported providers
	GetSupportedProviders() []types.VideoProvider
	
	// GetSupportedLanguages returns list of supported translation languages
	GetSupportedLanguages() []types.Language
}

// ProviderServiceInterface defines the contract for individual provider services
type ProviderServiceInterface interface {
	// GetVideoInfo retrieves basic information about a video
	GetVideoInfo(ctx context.Context, videoID string) (*types.VideoInfo, error)
	
	// GetTranscript retrieves transcript for a video in specified language
	GetTranscript(ctx context.Context, videoID string, language string) (*types.Transcript, error)
	
	// GetAvailableLanguages returns list of available transcript languages
	GetAvailableLanguages(ctx context.Context, videoID string) ([]types.Language, error)
	
	// GetCapabilities returns what features are supported for this video
	GetCapabilities(ctx context.Context, videoID string) (*types.VideoCapabilities, error)
	
	// GetProvider returns the video provider this service handles
	GetProvider() types.VideoProvider
	
	// ValidateVideoID checks if the video ID is valid for this provider
	ValidateVideoID(videoID string) bool
}

// ProviderFactory creates video service instances based on provider
type ProviderFactory interface {
	CreateService(provider types.VideoProvider) (ProviderServiceInterface, error)
	GetSupportedProviders() []types.VideoProvider
}