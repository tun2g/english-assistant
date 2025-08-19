package youtube

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app-backend/internal/types"
	oauthService "app-backend/internal/services/oauth"
	"go.uber.org/zap"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/option"
	"golang.org/x/oauth2"
)

// Service implements video.ServiceInterface for YouTube
type Service struct {
	apiKey      string
	service     *youtube.Service
	httpClient  *http.Client
	logger      *zap.Logger
	oauthService oauthService.ServiceInterface
}

// NewService creates a new YouTube service instance
func NewService(apiKey string, logger *zap.Logger) *Service {
	return NewServiceWithOAuth(apiKey, nil, logger)
}

// NewServiceWithOAuth creates a new YouTube service instance with OAuth support
func NewServiceWithOAuth(apiKey string, oauthSvc oauthService.ServiceInterface, logger *zap.Logger) *Service {
	ytService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		logger.Error("Failed to create youtube service", zap.Error(err))
		return &Service{
			apiKey:       apiKey,
			service:      nil, // Will cause graceful degradation
			httpClient:   &http.Client{Timeout: 30 * time.Second},
			logger:       logger,
			oauthService: oauthSvc,
		}
	}

	return &Service{
		apiKey:       apiKey,
		service:      ytService,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		logger:       logger,
		oauthService: oauthSvc,
	}
}

// GetProvider returns the YouTube provider identifier
func (s *Service) GetProvider() types.VideoProvider {
	return types.ProviderYouTube
}

// ValidateVideoID checks if the video ID is a valid YouTube video ID
func (s *Service) ValidateVideoID(videoID string) bool {
	// YouTube video IDs are 11 characters long and contain alphanumeric characters, hyphens, and underscores
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]{11}$`, videoID)
	return matched
}

// GetVideoInfo retrieves basic information about a YouTube video
func (s *Service) GetVideoInfo(ctx context.Context, videoID string) (*types.VideoInfo, error) {
	if !s.ValidateVideoID(videoID) {
		return nil, fmt.Errorf("invalid YouTube video ID: %s", videoID)
	}

	call := s.service.Videos.List([]string{"snippet", "contentDetails"}).Id(videoID)
	response, err := call.Context(ctx).Do()
	if err != nil {
		s.logger.Error("Failed to get video info", zap.String("videoID", videoID), zap.Error(err))
		return nil, fmt.Errorf("failed to get video info: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found: %s", videoID)
	}

	video := response.Items[0]
	duration, _ := parseISO8601Duration(video.ContentDetails.Duration)

	return &types.VideoInfo{
		ID:          videoID,
		Provider:    types.ProviderYouTube,
		Title:       video.Snippet.Title,
		Description: video.Snippet.Description,
		Duration:    types.MillisecondDuration(duration),
		ThumbnailURL: video.Snippet.Thumbnails.High.Url,
		URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	}, nil
}

// GetTranscript retrieves transcript for a YouTube video
func (s *Service) GetTranscript(ctx context.Context, videoID string, language string) (*types.Transcript, error) {
	if !s.ValidateVideoID(videoID) {
		return nil, fmt.Errorf("invalid YouTube video ID: %s", videoID)
	}

	// First, get available captions
	captionsCall := s.service.Captions.List([]string{"snippet"}, videoID)
	captionsResponse, err := captionsCall.Context(ctx).Do()
	if err != nil {
		s.logger.Error("Failed to list captions", zap.String("videoID", videoID), zap.Error(err))
		return nil, fmt.Errorf("failed to list captions: %w", err)
	}

	if len(captionsResponse.Items) == 0 {
		return &types.Transcript{
			VideoID:   videoID,
			Provider:  types.ProviderYouTube,
			Available: false,
		}, nil
	}

	// Find the best caption track
	var selectedCaption *youtube.Caption
	for _, caption := range captionsResponse.Items {
		if language != "" && caption.Snippet.Language == language {
			selectedCaption = caption
			break
		}
		// Fallback to first available caption if no language specified
		if selectedCaption == nil {
			selectedCaption = caption
		}
	}

	if selectedCaption == nil {
		return &types.Transcript{
			VideoID:   videoID,
			Provider:  types.ProviderYouTube,
			Available: false,
		}, nil
	}

	// Download the caption via API first
	segments, err := s.downloadCaption(ctx, selectedCaption.Id)
	if err != nil {
		s.logger.Warn("API caption download failed, trying web scraping fallback", 
			zap.String("videoID", videoID), 
			zap.String("captionID", selectedCaption.Id), 
			zap.Error(err))
		
		// Try web scraping fallback when API fails (especially for 403 errors)
		segments, err = s.scrapeTranscript(ctx, videoID, language)
		if err != nil {
			s.logger.Error("Both API and scraping methods failed", zap.String("videoID", videoID), zap.Error(err))
			return nil, fmt.Errorf("failed to retrieve transcript: %w", err)
		}
		
		s.logger.Info("Successfully retrieved transcript via web scraping", zap.String("videoID", videoID))
	}

	return &types.Transcript{
		VideoID:   videoID,
		Provider:  types.ProviderYouTube,
		Language:  selectedCaption.Snippet.Language,
		Segments:  segments,
		Available: true,
		Source:    getTrackKind(selectedCaption.Snippet.TrackKind),
	}, nil
}

// GetAvailableLanguages returns list of available transcript languages
func (s *Service) GetAvailableLanguages(ctx context.Context, videoID string) ([]types.Language, error) {
	if !s.ValidateVideoID(videoID) {
		return nil, fmt.Errorf("invalid YouTube video ID: %s", videoID)
	}

	call := s.service.Captions.List([]string{"snippet"}, videoID)
	response, err := call.Context(ctx).Do()
	if err != nil {
		s.logger.Error("Failed to list captions", zap.String("videoID", videoID), zap.Error(err))
		return nil, fmt.Errorf("failed to list captions: %w", err)
	}

	var languages []types.Language
	for _, caption := range response.Items {
		languages = append(languages, types.Language{
			Code: caption.Snippet.Language,
			Name: caption.Snippet.Name,
		})
	}

	return languages, nil
}

// GetCapabilities returns what features are supported for this video
func (s *Service) GetCapabilities(ctx context.Context, videoID string) (*types.VideoCapabilities, error) {
	languages, err := s.GetAvailableLanguages(ctx, videoID)
	if err != nil {
		return nil, err
	}

	hasAutoGenerated := false
	call := s.service.Captions.List([]string{"snippet"}, videoID)
	response, err := call.Context(ctx).Do()
	if err == nil {
		for _, caption := range response.Items {
			if caption.Snippet.TrackKind == "asr" {
				hasAutoGenerated = true
				break
			}
		}
	}

	return &types.VideoCapabilities{
		HasTranscript:        len(languages) > 0,
		AvailableLanguages:   languages,
		SupportsAutoGenerated: hasAutoGenerated,
	}, nil
}

// downloadCaption downloads and parses the caption content using OAuth2
func (s *Service) downloadCaption(ctx context.Context, captionID string) ([]types.TranscriptSegment, error) {
	// Check if OAuth service is available
	if s.oauthService == nil {
		return nil, fmt.Errorf("OAuth service not available - YouTube Caption API requires authentication")
	}

	if !s.oauthService.IsAuthenticated() {
		return nil, fmt.Errorf("user not authenticated - please authenticate with YouTube to access captions")
	}

	// Get valid OAuth token
	token, err := s.oauthService.GetValidToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get valid OAuth token: %w", err)
	}

	// Download the actual caption using authenticated API call
	return s.downloadCaptionWithAuth(ctx, captionID, token)
}

// downloadCaptionWithAuth downloads caption using OAuth2 authentication
func (s *Service) downloadCaptionWithAuth(ctx context.Context, captionID string, token *oauth2.Token) ([]types.TranscriptSegment, error) {
	// Create authenticated HTTP client
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	
	// Use authenticated client to create YouTube service
	authService, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		s.logger.Error("Failed to create authenticated YouTube service", zap.Error(err))
		return nil, fmt.Errorf("failed to create authenticated YouTube service: %w", err)
	}

	// Download caption content
	call := authService.Captions.Download(captionID)
	resp, err := call.Context(ctx).Download()
	if err != nil {
		s.logger.Error("Failed to download caption", zap.String("captionID", captionID), zap.Error(err))
		return nil, fmt.Errorf("failed to download caption: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read caption response", zap.Error(err))
		return nil, fmt.Errorf("failed to read caption response: %w", err)
	}

	s.logger.Info("Successfully downloaded caption content", 
		zap.String("captionID", captionID), 
		zap.Int("bodySize", len(body)))

	// Parse the caption content (YouTube returns TTML format)
	return s.parseTTML(body)
}

// parseTTML parses TTML caption format from YouTube
func (s *Service) parseTTML(data []byte) ([]types.TranscriptSegment, error) {
	// First try to parse as XML (TTML format)
	var ttml TTMLDocument
	if err := xml.Unmarshal(data, &ttml); err != nil {
		// If XML parsing fails, try parsing as plain text (some captions are just text)
		s.logger.Debug("Failed to parse as TTML, trying plain text", zap.Error(err))
		return s.parsePlainTextCaption(string(data)), nil
	}

	var segments []types.TranscriptSegment
	segmentIndex := 0
	
	for _, p := range ttml.Body.Div.P {
		startTime, err := s.parseTimeCode(p.Begin)
		if err != nil {
			s.logger.Warn("Failed to parse start time", zap.String("time", p.Begin), zap.Error(err))
			continue
		}

		endTime, err := s.parseTimeCode(p.End)
		if err != nil {
			s.logger.Warn("Failed to parse end time", zap.String("time", p.End), zap.Error(err))
			continue
		}

		// Clean up the text (remove XML tags)
		text := s.cleanCaptionText(p.Text)
		if text != "" && len(strings.TrimSpace(text)) > 0 {
			segments = append(segments, types.TranscriptSegment{
				Text:      text,
				StartTime: types.MillisecondDuration(startTime.Nanoseconds() / 1000000),
				EndTime:   types.MillisecondDuration(endTime.Nanoseconds() / 1000000),
				Index:     segmentIndex,
			})
			segmentIndex++
		}
	}
	
	// If no segments found, the TTML structure might be different
	if len(segments) == 0 {
		s.logger.Warn("No segments found in TTML, trying alternative parsing")
		return s.parsePlainTextCaption(string(data)), nil
	}

	s.logger.Info("Successfully parsed TTML captions", zap.Int("segments", len(segments)))
	return segments, nil
}

// parsePlainTextCaption parses plain text captions as fallback
func (s *Service) parsePlainTextCaption(text string) []types.TranscriptSegment {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	var segments []types.TranscriptSegment
	segmentIndex := 0

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && len(line) > 0 {
			// Create segments with estimated timing
			startMs := int64(i * 3000) // 3 seconds per segment
			endMs := startMs + 3000
			
			segments = append(segments, types.TranscriptSegment{
				Text:      line,
				StartTime: types.MillisecondDuration(startMs),
				EndTime:   types.MillisecondDuration(endMs),
				Index:     segmentIndex,
			})
			segmentIndex++
		}
	}

	s.logger.Debug("Parsed plain text captions", zap.Int("segments", len(segments)))
	return segments
}

// parseTimeCode parses TTML time codes (e.g., "00:00:01.500", "1.5s")
func (s *Service) parseTimeCode(timeStr string) (time.Duration, error) {
	// Handle different time formats that YouTube uses
	if strings.Contains(timeStr, ":") {
		// Format: HH:MM:SS.mmm or MM:SS.mmm
		parts := strings.Split(timeStr, ":")
		if len(parts) < 2 {
			return 0, fmt.Errorf("invalid time format: %s", timeStr)
		}

		var hours, minutes int
		var secondsStr string

		if len(parts) == 3 {
			// HH:MM:SS.mmm
			var err error
			hours, err = strconv.Atoi(parts[0])
			if err != nil {
				return 0, err
			}
			minutes, err = strconv.Atoi(parts[1])
			if err != nil {
				return 0, err
			}
			secondsStr = parts[2]
		} else {
			// MM:SS.mmm
			var err error
			minutes, err = strconv.Atoi(parts[0])
			if err != nil {
				return 0, err
			}
			secondsStr = parts[1]
		}

		seconds, err := strconv.ParseFloat(secondsStr, 64)
		if err != nil {
			return 0, err
		}

		totalSeconds := float64(hours*3600 + minutes*60) + seconds
		return time.Duration(totalSeconds * float64(time.Second)), nil
	} else if strings.HasSuffix(timeStr, "s") {
		// Format: "1.5s"
		secondsStr := strings.TrimSuffix(timeStr, "s")
		seconds, err := strconv.ParseFloat(secondsStr, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(seconds * float64(time.Second)), nil
	}

	return 0, fmt.Errorf("unsupported time format: %s", timeStr)
}

// cleanCaptionText removes XML tags and cleans up caption text
func (s *Service) cleanCaptionText(text string) string {
	// Remove XML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text = re.ReplaceAllString(text, "")
	
	// Clean up whitespace
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	return text
}

// parseISO8601Duration parses YouTube's ISO 8601 duration format
func parseISO8601Duration(duration string) (time.Duration, error) {
	// YouTube uses ISO 8601 duration format: PT4M13S
	re := regexp.MustCompile(`PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?`)
	matches := re.FindStringSubmatch(duration)
	
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", duration)
	}

	var hours, minutes, seconds int
	var err error

	if matches[1] != "" {
		hours, err = strconv.Atoi(matches[1])
		if err != nil {
			return 0, err
		}
	}
	if matches[2] != "" {
		minutes, err = strconv.Atoi(matches[2])
		if err != nil {
			return 0, err
		}
	}
	if matches[3] != "" {
		seconds, err = strconv.Atoi(matches[3])
		if err != nil {
			return 0, err
		}
	}

	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}

// scrapeTranscript scrapes transcript data from YouTube's web interface
// This is a fallback when the official API fails due to permissions
func (s *Service) scrapeTranscript(ctx context.Context, videoID, language string) ([]types.TranscriptSegment, error) {
	s.logger.Info("Starting transcript scraping", zap.String("videoID", videoID), zap.String("language", language))
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// First, get the video page to extract transcript data
	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	req, err := http.NewRequestWithContext(ctx, "GET", videoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers to mimic a browser request (improved for better success)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video page: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch video page, status: %d", resp.StatusCode)
	}
	
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Extract transcript data from the page
	segments, err := s.extractTranscriptFromHTML(string(body), language)
	if err != nil {
		return nil, fmt.Errorf("failed to extract transcript from HTML: %w", err)
	}
	
	s.logger.Info("Successfully scraped transcript", 
		zap.String("videoID", videoID), 
		zap.Int("segments", len(segments)))
	
	return segments, nil
}

// extractTranscriptFromHTML extracts transcript data from YouTube's HTML page with improved patterns
func (s *Service) extractTranscriptFromHTML(html, language string) ([]types.TranscriptSegment, error) {
	s.logger.Debug("Attempting to extract transcript from HTML", zap.Int("htmlLength", len(html)))
	
	// Check if we have any caption-related content
	if strings.Contains(html, "captionTracks") {
		s.logger.Debug("Found captionTracks in HTML")
	} else {
		s.logger.Warn("No captionTracks found in HTML - video may not have transcripts")
	}
	
	var transcriptURL string
	
	// Comprehensive patterns to find transcript URLs - inspired by youtube-transcript-api
	patterns := []string{
		// Look for baseUrl in any context containing timedtext
		`"baseUrl"\s*:\s*"([^"]*timedtext[^"]*)"`,
		// Look for timedtext URLs directly 
		`https://www\.youtube\.com/api/timedtext[^"'\s\)\]>]+`,
		// Look in caption tracks context
		`"captionTracks"[^}]*?"baseUrl"\s*:\s*"([^"]+)"`,
		// Look for any timedtext URL in quotes
		`"(https://[^"]*timedtext[^"]*)"`,
		`'(https://[^']*timedtext[^']*)'`,
		// Very broad search for any timedtext URLs
		`(https://[^\s"'<>]+timedtext[^\s"'<>]*)`,
		// Search in ytInitialPlayerResponse context
		`ytInitialPlayerResponse[^}]*?"baseUrl"\s*:\s*"([^"]*timedtext[^"]*)"`,
	}
	
	for i, pattern := range patterns {
		s.logger.Debug("Trying pattern", zap.Int("patternIndex", i))
		
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(html)
		
		if len(matches) >= 1 {
			// Take the URL from the appropriate capture group
			url := matches[0]
			if len(matches) >= 2 && matches[1] != "" && strings.Contains(matches[1], "timedtext") {
				url = matches[1]
			}
			
			if strings.Contains(url, "timedtext") && strings.Contains(url, "youtube") {
				s.logger.Debug("Pattern matched", zap.Int("patternIndex", i), zap.String("url", url[:min(len(url), 100)]))
				transcriptURL = url
				break
			}
		}
	}
	
	if transcriptURL == "" {
		s.logger.Error("No transcript URL found after all extraction methods")
		return nil, fmt.Errorf("no transcript URL found in page HTML - this may indicate the video has no available transcripts")
	}
	
	// Clean up and decode the URL
	transcriptURL = s.cleanTranscriptURL(transcriptURL)
	
	s.logger.Info("Attempting to fetch transcript", zap.String("url", transcriptURL[:min(len(transcriptURL), 100)]))
	
	// Fetch the actual transcript data
	return s.fetchTranscriptFromURL(transcriptURL)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// cleanTranscriptURL cleans and decodes the transcript URL
func (s *Service) cleanTranscriptURL(url string) string {
	// Decode common URL escaping
	url = strings.ReplaceAll(url, "\\u0026", "&")
	url = strings.ReplaceAll(url, "\\u003d", "=")
	url = strings.ReplaceAll(url, "\\u003c", "<")
	url = strings.ReplaceAll(url, "\\u003e", ">")
	url = strings.ReplaceAll(url, "\\/", "/")
	url = strings.ReplaceAll(url, "\\", "")
	
	return url
}

// fetchTranscriptFromURL fetches and parses transcript data from the YouTube transcript URL
func (s *Service) fetchTranscriptFromURL(url string) ([]types.TranscriptSegment, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transcript, status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read transcript response: %w", err)
	}
	
	s.logger.Debug("Raw transcript response", zap.Int("bodyLength", len(body)), zap.String("contentType", resp.Header.Get("Content-Type")))
	
	// Parse the XML transcript data
	return s.parseTranscriptXML(string(body))
}

// parseTranscriptXML parses XML transcript data from YouTube
func (s *Service) parseTranscriptXML(xmlData string) ([]types.TranscriptSegment, error) {
	s.logger.Debug("Parsing transcript XML", zap.Int("xmlLength", len(xmlData)))
	
	// Enhanced patterns for parsing transcript XML based on youtube-transcript-api
	patterns := []string{
		// Standard format: <text start="0.0" dur="1.5">Hello world</text>
		`<text start="([^"]+)" dur="([^"]+)"[^>]*>([^<]*)</text>`,
		// Alternative format with different attributes
		`<text start="([^"]+)" duration="([^"]+)"[^>]*>([^<]*)</text>`,
		// Format with t attribute instead of start
		`<text t="([^"]+)" d="([^"]+)"[^>]*>([^<]*)</text>`,
		// Format with milliseconds
		`<p t="([^"]+)" d="([^"]+)"[^>]*>([^<]*)</p>`,
		// More flexible patterns for different attribute orders
		`<text[^>]*start="([^"]+)"[^>]*dur="([^"]+)"[^>]*>([^<]*)</text>`,
		// WebVTT-style patterns
		`<c t="([^"]+)" d="([^"]+)"[^>]*>([^<]*)</c>`,
	}
	
	var segments []types.TranscriptSegment
	
	for i, pattern := range patterns {
		s.logger.Debug("Trying XML pattern", zap.Int("patternIndex", i))
		
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(xmlData, -1)
		
		if len(matches) > 0 {
			s.logger.Debug("XML pattern matched", zap.Int("patternIndex", i), zap.Int("matches", len(matches)))
			
			for _, match := range matches {
				if len(match) < 4 {
					continue
				}
				
				startTime, err := strconv.ParseFloat(match[1], 64)
				if err != nil {
					s.logger.Warn("Failed to parse start time", zap.String("time", match[1]))
					continue
				}
				
				duration, err := strconv.ParseFloat(match[2], 64)
				if err != nil {
					s.logger.Warn("Failed to parse duration", zap.String("duration", match[2]))
					continue
				}
				
				text := html.UnescapeString(match[3])
				text = s.cleanCaptionText(text)
				
				if text == "" {
					continue
				}
				
				// Convert to MillisecondDuration
				startTimeMs := types.MillisecondDuration(time.Duration(startTime * float64(time.Second)))
				endTimeMs := types.MillisecondDuration(time.Duration((startTime + duration) * float64(time.Second)))
				
				segments = append(segments, types.TranscriptSegment{
					Text:      text,
					StartTime: startTimeMs,
					EndTime:   endTimeMs,
					Index:     len(segments),
				})
			}
			
			if len(segments) > 0 {
				break // Successfully parsed with this pattern
			}
		}
	}
	
	if len(segments) == 0 {
		// Try parsing as plain text with timestamps if XML parsing fails
		return s.parseAsPlainText(xmlData)
	}
	
	s.logger.Info("Successfully parsed transcript XML", zap.Int("segments", len(segments)))
	return segments, nil
}

// parseAsPlainText tries to parse transcript as plain text if XML parsing fails
func (s *Service) parseAsPlainText(data string) ([]types.TranscriptSegment, error) {
	s.logger.Debug("Attempting to parse as plain text")
	
	// If the data looks like it contains timing info, try to extract it
	lines := strings.Split(data, "\n")
	var segments []types.TranscriptSegment
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Simple fallback: create segments with estimated timing
		// This is a basic approach when we can't parse proper timing
		startTime := float64(i * 3) // Assume 3 seconds per segment
		duration := 3.0
		
		startTimeMs := types.MillisecondDuration(time.Duration(startTime * float64(time.Second)))
		endTimeMs := types.MillisecondDuration(time.Duration((startTime + duration) * float64(time.Second)))
		
		segments = append(segments, types.TranscriptSegment{
			Text:      line,
			StartTime: startTimeMs,
			EndTime:   endTimeMs,
			Index:     len(segments),
		})
		
		// Limit to prevent too many segments from malformed data
		if len(segments) >= 100 {
			break
		}
	}
	
	if len(segments) == 0 {
		return nil, fmt.Errorf("no valid transcript segments extracted from any format")
	}
	
	s.logger.Info("Parsed transcript as plain text", zap.Int("segments", len(segments)))
	return segments, nil
}

// getTrackKind converts YouTube track kind to our format
func getTrackKind(trackKind string) string {
	switch trackKind {
	case "asr":
		return "auto-generated"
	case "forced":
		return "forced"
	default:
		return "manual"
	}
}