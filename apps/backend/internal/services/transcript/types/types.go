package types

import "time"

// TranscriptSegment represents a single segment of transcript text with timing
type TranscriptSegment struct {
	Text     string        `json:"text"`
	Start    time.Duration `json:"start"`
	Duration time.Duration `json:"duration"`
	Offset   int64         `json:"offset,omitempty"`
}

// Transcript represents the complete transcript of a video
type Transcript struct {
	VideoID    string               `json:"video_id"`
	Title      string               `json:"title,omitempty"`
	Language   string               `json:"language"`
	Segments   []TranscriptSegment  `json:"segments"`
	Provider   string               `json:"provider"`
	CreatedAt  time.Time            `json:"created_at"`
}

// TranscriptRequest represents a request for video transcript
type TranscriptRequest struct {
	VideoID     string `json:"video_id" validate:"required"`
	VideoURL    string `json:"video_url,omitempty"`
	Language    string `json:"language,omitempty"`
	Country     string `json:"country,omitempty"`
	PreferredProviders []string `json:"preferred_providers,omitempty"`
}

// ProviderType represents available transcript providers
type ProviderType string

const (
	ProviderYouTubeAPI    ProviderType = "youtube_api"
	ProviderYTTranscript  ProviderType = "yt_transcript"
	ProviderKkdaiYouTube  ProviderType = "kkdai_youtube"
	ProviderInnertube     ProviderType = "innertube"
)

// ProviderConfig represents configuration for a transcript provider
type ProviderConfig struct {
	Type     ProviderType `json:"type"`
	Enabled  bool         `json:"enabled"`
	Priority int          `json:"priority"`
	Config   map[string]interface{} `json:"config,omitempty"`
}