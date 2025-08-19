package youtube_api

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"

	"app-backend/internal/logger"
	"app-backend/internal/services/transcript/errors"
	"app-backend/internal/services/transcript/types"
)

type Provider struct {
	apiKey   string
	service  *youtube.Service
	logger   *logger.Logger
	priority int
}

type Config struct {
	APIKey   string `json:"api_key"`
	Priority int    `json:"priority"`
}

func NewProvider(config *Config, logger *logger.Logger) (*Provider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("YouTube API key is required")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %w", err)
	}

	priority := config.Priority
	if priority == 0 {
		priority = 1 // Default priority
	}

	return &Provider{
		apiKey:   config.APIKey,
		service:  service,
		logger:   logger,
		priority: priority,
	}, nil
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

	// Get video details
	videoCall := p.service.Videos.List([]string{"snippet"}).Id(videoID)
	videoResponse, err := videoCall.Do()
	if err != nil {
		return nil, errors.NewProviderError("youtube_api", err)
	}

	if len(videoResponse.Items) == 0 {
		return nil, errors.ErrTranscriptNotFound
	}

	video := videoResponse.Items[0]

	// List available captions
	captionsCall := p.service.Captions.List([]string{"snippet"}, videoID)
	captionsResponse, err := captionsCall.Do()
	if err != nil {
		return nil, errors.NewProviderError("youtube_api", err)
	}

	if len(captionsResponse.Items) == 0 {
		return nil, errors.ErrTranscriptNotFound
	}

	// Find the best caption track
	var selectedCaption *youtube.Caption
	language := req.Language
	if language == "" {
		language = "en"
	}

	// Try to find exact language match first
	for _, caption := range captionsResponse.Items {
		if caption.Snippet.Language == language {
			selectedCaption = caption
			break
		}
	}

	// If no exact match, try language prefix (e.g., "en" for "en-US")
	if selectedCaption == nil {
		languagePrefix := strings.Split(language, "-")[0]
		for _, caption := range captionsResponse.Items {
			if strings.HasPrefix(caption.Snippet.Language, languagePrefix) {
				selectedCaption = caption
				break
			}
		}
	}

	// If still no match, use first available caption
	if selectedCaption == nil {
		selectedCaption = captionsResponse.Items[0]
	}

	// Download caption content
	downloadCall := p.service.Captions.Download(selectedCaption.Id).Tfmt("srt")
	response, err := downloadCall.Download()
	if err != nil {
		return nil, errors.NewProviderError("youtube_api", err)
	}
	defer response.Body.Close()

	// Read response body
	buf := make([]byte, 1024*1024) // 1MB buffer
	n, err := response.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return nil, errors.NewProviderError("youtube_api", err)
	}

	srtContent := string(buf[:n])

	// Parse SRT content
	segments, err := p.parseSRT(srtContent)
	if err != nil {
		return nil, errors.NewProviderError("youtube_api", err)
	}

	return &types.Transcript{
		VideoID:   videoID,
		Title:     video.Snippet.Title,
		Language:  selectedCaption.Snippet.Language,
		Segments:  segments,
		Provider:  string(types.ProviderYouTubeAPI),
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
	// Test API availability with a simple call
	_, err := p.service.Videos.List([]string{"snippet"}).Id("dQw4w9WgXcQ").Do()
	return err == nil
}

func (p *Provider) GetProviderType() types.ProviderType {
	return types.ProviderYouTubeAPI
}

func (p *Provider) GetPriority() int {
	return p.priority
}

// parseSRT parses SRT subtitle format into transcript segments
func (p *Provider) parseSRT(content string) ([]types.TranscriptSegment, error) {
	var segments []types.TranscriptSegment
	
	blocks := strings.Split(content, "\n\n")
	for _, block := range blocks {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) < 3 {
			continue
		}

		// Parse timing line (format: 00:00:01,000 --> 00:00:04,000)
		timingLine := lines[1]
		times := strings.Split(timingLine, " --> ")
		if len(times) != 2 {
			continue
		}

		start, err := p.parseSRTTime(strings.TrimSpace(times[0]))
		if err != nil {
			continue
		}

		end, err := p.parseSRTTime(strings.TrimSpace(times[1]))
		if err != nil {
			continue
		}

		// Combine text lines
		text := strings.Join(lines[2:], " ")
		text = strings.TrimSpace(text)

		if text != "" {
			segments = append(segments, types.TranscriptSegment{
				Text:     text,
				Start:    start,
				Duration: end - start,
			})
		}
	}

	return segments, nil
}

// parseSRTTime parses SRT time format (00:00:01,000) to time.Duration
func (p *Provider) parseSRTTime(timeStr string) (time.Duration, error) {
	// Replace comma with dot for milliseconds
	timeStr = strings.Replace(timeStr, ",", ".", 1)
	
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}

	secondsParts := strings.Split(parts[2], ".")
	seconds, err := strconv.Atoi(secondsParts[0])
	if err != nil {
		return 0, err
	}

	var milliseconds int
	if len(secondsParts) > 1 {
		// Pad or truncate to 3 digits
		msStr := secondsParts[1]
		if len(msStr) > 3 {
			msStr = msStr[:3]
		} else {
			for len(msStr) < 3 {
				msStr += "0"
			}
		}
		milliseconds, err = strconv.Atoi(msStr)
		if err != nil {
			return 0, err
		}
	}

	duration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(milliseconds)*time.Millisecond

	return duration, nil
}