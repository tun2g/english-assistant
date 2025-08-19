package yt_transcript

import (
	"context"
	"time"

	"github.com/chand1012/yt_transcript"
	"go.uber.org/zap"

	"app-backend/internal/logger"
	"app-backend/internal/services/transcript/errors"
	"app-backend/internal/services/transcript/types"
)

type Provider struct {
	logger   *logger.Logger
	priority int
}

type Config struct {
	Priority int `json:"priority"`
}

func NewProvider(config *Config, logger *logger.Logger) *Provider {
	priority := config.Priority
	if priority == 0 {
		priority = 2 // Default priority (lower than YouTube API)
	}

	return &Provider{
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

	language := req.Language
	if language == "" {
		language = "en"
	}

	country := req.Country
	if country == "" {
		country = "US"
	}

	// Fetch transcript using yt_transcript library
	transcriptResponses, title, err := yt_transcript.FetchTranscript(videoID, language, country)
	if err != nil {
		p.logger.Error("Failed to fetch transcript with yt_transcript", 
			zap.String("video_id", videoID),
			zap.String("language", language),
			zap.String("country", country),
			zap.Error(err))
		return nil, errors.NewProviderError("yt_transcript", err)
	}

	if len(transcriptResponses) == 0 {
		return nil, errors.ErrTranscriptNotFound
	}

	// Convert to our transcript format
	segments := make([]types.TranscriptSegment, len(transcriptResponses))
	for i, resp := range transcriptResponses {
		segments[i] = types.TranscriptSegment{
			Text:     resp.Text,
			Start:    time.Duration(resp.Offset) * time.Millisecond,
			Duration: time.Duration(resp.Duration) * time.Millisecond,
		}
	}

	return &types.Transcript{
		VideoID:   videoID,
		Title:     title,
		Language:  language,
		Segments:  segments,
		Provider:  string(types.ProviderYTTranscript),
		CreatedAt: time.Now(),
	}, nil
}

func (p *Provider) GetVideoID(url string) (string, error) {
	return yt_transcript.GetVideoID(url)
}

func (p *Provider) IsAvailable(ctx context.Context) bool {
	// Test with a known video that should have transcripts
	_, _, err := yt_transcript.FetchTranscript("dQw4w9WgXcQ", "en", "US")
	return err == nil
}

func (p *Provider) GetProviderType() types.ProviderType {
	return types.ProviderYTTranscript
}

func (p *Provider) GetPriority() int {
	return p.priority
}