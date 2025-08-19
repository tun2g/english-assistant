package transcript

import (
	"context"

	"app-backend/internal/services/transcript/types"
)

// ProviderInterface defines the contract for transcript providers
type ProviderInterface interface {
	// GetTranscript retrieves transcript for a video
	GetTranscript(ctx context.Context, req *types.TranscriptRequest) (*types.Transcript, error)
	
	// GetVideoID extracts video ID from YouTube URL
	GetVideoID(url string) (string, error)
	
	// IsAvailable checks if the provider is currently available
	IsAvailable(ctx context.Context) bool
	
	// GetProviderType returns the provider type
	GetProviderType() types.ProviderType
	
	// GetPriority returns the provider priority (lower number = higher priority)
	GetPriority() int
}

// ServiceInterface defines the main transcript service contract
type ServiceInterface interface {
	// GetTranscript retrieves transcript using the best available provider
	GetTranscript(ctx context.Context, req *types.TranscriptRequest) (*types.Transcript, error)
	
	// GetTranscriptWithProvider retrieves transcript using a specific provider
	GetTranscriptWithProvider(ctx context.Context, provider types.ProviderType, req *types.TranscriptRequest) (*types.Transcript, error)
	
	// GetAvailableProviders returns list of currently available providers
	GetAvailableProviders(ctx context.Context) []types.ProviderType
	
	// RegisterProvider adds a new provider to the service
	RegisterProvider(provider ProviderInterface) error
}