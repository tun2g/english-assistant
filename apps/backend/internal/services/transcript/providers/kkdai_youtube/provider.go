package kkdai_youtube

import (
	"context"
	"regexp"
	"time"

	"github.com/kkdai/youtube/v2"
	"go.uber.org/zap"

	"app-backend/internal/logger"
	"app-backend/internal/services/transcript/errors"
	"app-backend/internal/services/transcript/types"
)

type Provider struct {
	client   *youtube.Client
	logger   *logger.Logger
	priority int
}

type Config struct {
	Priority int `json:"priority"`
}

func NewProvider(config *Config, logger *logger.Logger) *Provider {
	priority := config.Priority
	if priority == 0 {
		priority = 3 // Default priority
	}

	return &Provider{
		client:   &youtube.Client{},
		logger:   logger,
		priority: priority,
	}
}

func (p *Provider) GetTranscript(ctx context.Context, req *types.TranscriptRequest) (*types.Transcript, error) {
	videoID := req.VideoID
	if videoID == "" && req.VideoURL != "" {
		var err error
		videoID, err = p.GetVideoID(req.VideoURL)
		if err != nil {
			return nil, err
		}
	}

	if videoID == "" {
		return nil, errors.ErrInvalidVideoID
	}

	// Get video information
	video, err := p.client.GetVideo(videoID)
	if err != nil {
		p.logger.Error("Failed to get video with kkdai/youtube", 
			zap.String("video_id", videoID),
			zap.Error(err))
		return nil, errors.NewProviderError("kkdai_youtube", err)
	}

	// Determine language from request or default to English
	language := req.Language
	if language == "" {
		language = "en"
	}

	// Get transcript
	transcript, err := p.client.GetTranscript(video, language)
	if err != nil {
		p.logger.Error("Failed to get transcript with kkdai/youtube", 
			zap.String("video_id", videoID),
			zap.String("language", language),
			zap.Error(err))
		
		// Check if it's the specific "transcript disabled" error
		if err == youtube.ErrTranscriptDisabled {
			return nil, errors.ErrTranscriptDisabled
		}
		
		return nil, errors.NewProviderError("kkdai_youtube", err)
	}

	if len(transcript) == 0 {
		return nil, errors.ErrTranscriptNotFound
	}

	// Convert to our transcript format
	segments := make([]types.TranscriptSegment, len(transcript))
	for i, segment := range transcript {
		segments[i] = types.TranscriptSegment{
			Text:     segment.Text,
			Start:    time.Duration(segment.StartMs) * time.Millisecond,
			Duration: time.Duration(segment.Duration) * time.Millisecond,
		}
	}

	return &types.Transcript{
		VideoID:   videoID,
		Title:     video.Title,
		Language:  language,
		Segments:  segments,
		Provider:  string(types.ProviderKkdaiYouTube),
		CreatedAt: time.Now(),
	}, nil
}

func (p *Provider) GetVideoID(url string) (string, error) {
	patterns := []string{
		`(?:youtube\.com/watch\?v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`,
		`(?:youtube\.com/v/)([a-zA-Z0-9_-]{11})`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Check if it's already a video ID
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{11}$`, url); matched {
		return url, nil
	}

	return "", errors.NewVideoIDExtractionError(url, nil)
}

func (p *Provider) IsAvailable(ctx context.Context) bool {
	// Test with a known video that should be accessible
	_, err := p.client.GetVideo("dQw4w9WgXcQ")
	return err == nil
}

func (p *Provider) GetProviderType() types.ProviderType {
	return types.ProviderKkdaiYouTube
}

func (p *Provider) GetPriority() int {
	return p.priority
}