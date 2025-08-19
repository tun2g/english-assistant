package transcript

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"go.uber.org/zap"

	"app-backend/internal/config"
	"app-backend/internal/logger"
	"app-backend/internal/services/transcript/errors"
	"app-backend/internal/services/transcript/providers/innertube"
	"app-backend/internal/services/transcript/providers/kkdai_youtube"
	"app-backend/internal/services/transcript/providers/yt_transcript"
	"app-backend/internal/services/transcript/providers/youtube_api"
	"app-backend/internal/services/transcript/types"
)

type Service struct {
	providers map[types.ProviderType]ProviderInterface
	config    *config.Config
	logger    *logger.Logger
	mu        sync.RWMutex
}

func NewService(config *config.Config, logger *logger.Logger) (*Service, error) {
	service := &Service{
		providers: make(map[types.ProviderType]ProviderInterface),
		config:    config,
		logger:    logger,
	}

	// Initialize providers based on configuration
	if err := service.initializeProviders(); err != nil {
		return nil, fmt.Errorf("failed to initialize providers: %w", err)
	}

	return service, nil
}

func (s *Service) initializeProviders() error {
	// Initialize YouTube API provider if configured
	if s.config.ExternalAPIs.YouTube.APIKey != "" {
		youtubeConfig := &youtube_api.Config{
			APIKey:   s.config.ExternalAPIs.YouTube.APIKey,
			Priority: 1,
		}
		provider, err := youtube_api.NewProvider(youtubeConfig, s.logger)
		if err != nil {
			s.logger.Warn("Failed to initialize YouTube API provider", zap.Error(err))
		} else {
			s.providers[types.ProviderYouTubeAPI] = provider
		}
	}

	// Initialize yt_transcript provider
	ytTranscriptConfig := &yt_transcript.Config{
		Priority: 2,
	}
	ytTranscriptProvider := yt_transcript.NewProvider(ytTranscriptConfig, s.logger)
	s.providers[types.ProviderYTTranscript] = ytTranscriptProvider

	// Initialize kkdai/youtube provider
	kkdaiConfig := &kkdai_youtube.Config{
		Priority: 3,
	}
	kkdaiProvider := kkdai_youtube.NewProvider(kkdaiConfig, s.logger)
	s.providers[types.ProviderKkdaiYouTube] = kkdaiProvider

	// Initialize Innertube provider
	innertubeConfig := &innertube.Config{
		Priority: 4,
		Timeout:  30,
	}
	innertubeProvider := innertube.NewProvider(innertubeConfig, s.logger)
	s.providers[types.ProviderInnertube] = innertubeProvider

	s.logger.Info("Initialized transcript providers", 
		zap.Int("provider_count", len(s.providers)),
		zap.Strings("providers", s.getProviderTypes()))

	return nil
}

func (s *Service) GetTranscript(ctx context.Context, req *types.TranscriptRequest) (*types.Transcript, error) {
	if req == nil {
		return nil, fmt.Errorf("transcript request cannot be nil")
	}

	// Validate request
	if req.VideoID == "" && req.VideoURL == "" {
		return nil, errors.ErrInvalidVideoID
	}

	// Get providers in priority order
	providers := s.getProvidersInPriorityOrder(req.PreferredProviders)
	if len(providers) == 0 {
		return nil, errors.ErrProviderNotAvailable
	}

	var lastErr error
	var providerErrors []string
	
	for _, provider := range providers {
		s.logger.Info("Attempting to get transcript", 
			zap.String("provider", string(provider.GetProviderType())),
			zap.String("video_id", req.VideoID),
			zap.String("video_url", req.VideoURL),
			zap.String("language", req.Language))

		// Check if provider is available
		if !provider.IsAvailable(ctx) {
			errMsg := fmt.Sprintf("Provider %s not available", provider.GetProviderType())
			providerErrors = append(providerErrors, errMsg)
			s.logger.Warn("Provider not available", 
				zap.String("provider", string(provider.GetProviderType())))
			continue
		}

		transcript, err := provider.GetTranscript(ctx, req)
		if err != nil {
			errMsg := fmt.Sprintf("Provider %s failed: %v", provider.GetProviderType(), err)
			providerErrors = append(providerErrors, errMsg)
			s.logger.Error("Provider failed to get transcript", 
				zap.String("provider", string(provider.GetProviderType())),
				zap.String("video_id", req.VideoID),
				zap.Error(err))
			lastErr = err
			continue
		}

		s.logger.Info("Successfully retrieved transcript", 
			zap.String("provider", string(provider.GetProviderType())),
			zap.String("video_id", transcript.VideoID),
			zap.Int("segment_count", len(transcript.Segments)),
			zap.String("language", transcript.Language))

		return transcript, nil
	}

	// Log summary of all failures
	s.logger.Error("All transcript providers failed", 
		zap.String("video_id", req.VideoID),
		zap.Strings("provider_errors", providerErrors),
		zap.Int("total_providers", len(providers)))

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, errors.ErrAllProvidersFailed
}

func (s *Service) GetTranscriptWithProvider(ctx context.Context, providerType types.ProviderType, req *types.TranscriptRequest) (*types.Transcript, error) {
	s.mu.RLock()
	provider, exists := s.providers[providerType]
	s.mu.RUnlock()

	if !exists {
		return nil, errors.ErrProviderNotAvailable
	}

	if !provider.IsAvailable(ctx) {
		return nil, errors.ErrProviderNotAvailable
	}

	return provider.GetTranscript(ctx, req)
}

func (s *Service) GetAvailableProviders(ctx context.Context) []types.ProviderType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var available []types.ProviderType
	for providerType, provider := range s.providers {
		if provider.IsAvailable(ctx) {
			available = append(available, providerType)
		}
	}

	return available
}

func (s *Service) RegisterProvider(provider ProviderInterface) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	providerType := provider.GetProviderType()
	s.providers[providerType] = provider

	s.logger.Info("Registered new transcript provider", 
		zap.String("provider", string(providerType)),
		zap.Int("priority", provider.GetPriority()))

	return nil
}

// getProvidersInPriorityOrder returns providers sorted by priority
// If preferred providers are specified, they are tried first in the order given
func (s *Service) getProvidersInPriorityOrder(preferredProviders []string) []ProviderInterface {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []ProviderInterface
	usedProviders := make(map[types.ProviderType]bool)

	// First, add preferred providers in the order specified
	for _, preferred := range preferredProviders {
		providerType := types.ProviderType(preferred)
		if provider, exists := s.providers[providerType]; exists {
			result = append(result, provider)
			usedProviders[providerType] = true
		}
	}

	// Then add remaining providers sorted by priority
	var remaining []ProviderInterface
	for providerType, provider := range s.providers {
		if !usedProviders[providerType] {
			remaining = append(remaining, provider)
		}
	}

	// Sort remaining providers by priority (lower number = higher priority)
	sort.Slice(remaining, func(i, j int) bool {
		return remaining[i].GetPriority() < remaining[j].GetPriority()
	})

	result = append(result, remaining...)
	return result
}

func (s *Service) getProviderTypes() []string {
	var types []string
	for providerType := range s.providers {
		types = append(types, string(providerType))
	}
	return types
}

// Health check methods
func (s *Service) HealthCheck(ctx context.Context) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := make(map[string]interface{})
	status["total_providers"] = len(s.providers)

	providerStatus := make(map[string]bool)
	availableCount := 0

	for providerType, provider := range s.providers {
		isAvailable := provider.IsAvailable(ctx)
		providerStatus[string(providerType)] = isAvailable
		if isAvailable {
			availableCount++
		}
	}

	status["available_providers"] = availableCount
	status["provider_status"] = providerStatus
	status["healthy"] = availableCount > 0

	return status
}