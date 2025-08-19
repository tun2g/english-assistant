package innertube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"app-backend/internal/logger"
	"app-backend/internal/services/transcript/errors"
	"app-backend/internal/services/transcript/types"
)

type Provider struct {
	httpClient *http.Client
	logger     *logger.Logger
	priority   int
}

type Config struct {
	Priority int `json:"priority"`
	Timeout  int `json:"timeout"` // in seconds
}

// Innertube API request structures
type InnertubeRequest struct {
	Context struct {
		Client struct {
			ClientName    string `json:"clientName"`
			ClientVersion string `json:"clientVersion"`
			Platform      string `json:"platform"`
		} `json:"client"`
	} `json:"context"`
	VideoID string `json:"videoId"`
}

// Innertube API response structures
type InnertubeResponse struct {
	Actions []struct {
		UpdateEngagementPanelAction struct {
			Content struct {
				TranscriptRenderer struct {
					Body struct {
						TranscriptBodyRenderer struct {
							CueGroups []struct {
								TranscriptCueGroupRenderer struct {
									Cues []struct {
										TranscriptCueRenderer struct {
											Cue struct {
												SimpleText string `json:"simpleText"`
											} `json:"cue"`
											StartOffsetMs string `json:"startOffsetMs"`
											DurationMs    string `json:"durationMs"`
										} `json:"transcriptCueRenderer"`
									} `json:"cues"`
								} `json:"transcriptCueGroupRenderer"`
							} `json:"cueGroups"`
						} `json:"transcriptBodyRenderer"`
					} `json:"body"`
				} `json:"transcriptRenderer"`
			} `json:"content"`
		} `json:"updateEngagementPanelAction"`
	} `json:"actions"`
}

func NewProvider(config *Config, logger *logger.Logger) *Provider {
	priority := config.Priority
	if priority == 0 {
		priority = 4 // Default priority
	}

	timeout := time.Duration(config.Timeout) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Provider{
		httpClient: &http.Client{
			Timeout: timeout,
		},
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

	// First, get video info to get title
	title, err := p.getVideoTitle(ctx, videoID)
	if err != nil {
		p.logger.Warn("Failed to get video title", 
			zap.String("video_id", videoID),
			zap.Error(err))
		title = "" // Continue without title
	}

	// Get transcript using Innertube API
	segments, language, err := p.fetchTranscriptFromInnertube(ctx, videoID, req.Language)
	if err != nil {
		return nil, err
	}

	if len(segments) == 0 {
		return nil, errors.ErrTranscriptNotFound
	}

	return &types.Transcript{
		VideoID:   videoID,
		Title:     title,
		Language:  language,
		Segments:  segments,
		Provider:  string(types.ProviderInnertube),
		CreatedAt: time.Now(),
	}, nil
}

func (p *Provider) fetchTranscriptFromInnertube(ctx context.Context, videoID, preferredLanguage string) ([]types.TranscriptSegment, string, error) {
	// Create Innertube request (Android client for better compatibility)
	innertubeReq := InnertubeRequest{
		VideoID: videoID,
	}
	innertubeReq.Context.Client.ClientName = "ANDROID"
	innertubeReq.Context.Client.ClientVersion = "17.31.35"
	innertubeReq.Context.Client.Platform = "MOBILE"

	reqBody, err := json.Marshal(innertubeReq)
	if err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}

	// Make request to Innertube API
	url := "https://www.youtube.com/youtubei/v1/get_transcript?key=AIzaSyA8eiZmM1FaDVjRy-df2KTyQ_vz_yYM39w"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "com.google.android.youtube/17.31.35 (Linux; U; Android 11) gzip")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		p.logger.Error("Innertube API error", 
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
			zap.String("video_id", videoID))
		return nil, "", errors.NewProviderError("innertube", fmt.Errorf("HTTP %d", resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}

	// Parse response
	var innertubeResp InnertubeResponse
	if err := json.Unmarshal(body, &innertubeResp); err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}

	// Extract transcript segments
	segments, err := p.parseInnertubeResponse(&innertubeResp)
	if err != nil {
		return nil, "", errors.NewProviderError("innertube", err)
	}

	language := preferredLanguage
	if language == "" {
		language = "en" // Default to English
	}

	return segments, language, nil
}

func (p *Provider) parseInnertubeResponse(resp *InnertubeResponse) ([]types.TranscriptSegment, error) {
	var segments []types.TranscriptSegment

	for _, action := range resp.Actions {
		transcriptRenderer := action.UpdateEngagementPanelAction.Content.TranscriptRenderer
		bodyRenderer := transcriptRenderer.Body.TranscriptBodyRenderer

		for _, cueGroup := range bodyRenderer.CueGroups {
			for _, cue := range cueGroup.TranscriptCueGroupRenderer.Cues {
				cueRenderer := cue.TranscriptCueRenderer

				text := cueRenderer.Cue.SimpleText
				if text == "" {
					continue
				}

				// Parse timing
				startMs, err := strconv.ParseInt(cueRenderer.StartOffsetMs, 10, 64)
				if err != nil {
					continue
				}

				durationMs, err := strconv.ParseInt(cueRenderer.DurationMs, 10, 64)
				if err != nil {
					continue
				}

				segment := types.TranscriptSegment{
					Text:     strings.TrimSpace(text),
					Start:    time.Duration(startMs) * time.Millisecond,
					Duration: time.Duration(durationMs) * time.Millisecond,
				}

				segments = append(segments, segment)
			}
		}
	}

	return segments, nil
}

func (p *Provider) getVideoTitle(ctx context.Context, videoID string) (string, error) {
	// Use a simple approach to get video title from YouTube page
	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract title using regex
	titleRegex := regexp.MustCompile(`<title>(.+?) - YouTube</title>`)
	matches := titleRegex.FindSubmatch(body)
	if len(matches) > 1 {
		return string(matches[1]), nil
	}

	// Alternative regex for different title formats
	titleRegex2 := regexp.MustCompile(`"title":"([^"]+)"`)
	matches2 := titleRegex2.FindSubmatch(body)
	if len(matches2) > 1 {
		return string(matches2[1]), nil
	}

	return "", fmt.Errorf("title not found")
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
	// Test with a simple request to YouTube
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.youtube.com", nil)
	if err != nil {
		return false
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (p *Provider) GetProviderType() types.ProviderType {
	return types.ProviderInnertube
}

func (p *Provider) GetPriority() int {
	return p.priority
}