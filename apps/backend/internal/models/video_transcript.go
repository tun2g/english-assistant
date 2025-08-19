package models

import (
	"time"
	"app-backend/internal/types"
)

// VideoTranscriptCache represents cached transcript data
type VideoTranscriptCache struct {
	Auditable
	VideoID      string                `gorm:"index;not null" json:"videoId"`
	Provider     types.VideoProvider   `gorm:"not null" json:"provider"`
	Language     string                `gorm:"not null" json:"language"`
	Content      string                `gorm:"type:text;not null" json:"content"` // JSON-encoded transcript segments
	Source       string                `gorm:"default:'manual'" json:"source"`    // "manual", "auto-generated"
	ExpiresAt    time.Time             `gorm:"index" json:"expiresAt"`
}

// VideoTranslationCache represents cached translation data
type VideoTranslationCache struct {
	Auditable
	VideoID      string                `gorm:"index;not null" json:"videoId"`
	Provider     types.VideoProvider   `gorm:"not null" json:"provider"`
	SourceLang   string                `gorm:"not null" json:"sourceLang"`
	TargetLang   string                `gorm:"not null" json:"targetLang"`
	Content      string                `gorm:"type:text;not null" json:"content"` // JSON-encoded translated segments
	ExpiresAt    time.Time             `gorm:"index" json:"expiresAt"`
}

// UserAPIKey represents encrypted API keys for users
type UserAPIKey struct {
	Auditable
	UserID       uint   `gorm:"index;not null" json:"userId"`
	ServiceName  string `gorm:"not null" json:"serviceName"` // "youtube", "gemini"
	EncryptedKey string `gorm:"not null" json:"-"`           // Don't expose in JSON
	KeyHash      string `gorm:"index" json:"-"`              // Hash for quick lookup
	
	// Relationship
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// VideoAnalytics represents usage analytics for videos
type VideoAnalytics struct {
	Auditable
	VideoID           string              `gorm:"index;not null" json:"videoId"`
	Provider          types.VideoProvider `gorm:"not null" json:"provider"`
	UserID            uint                `gorm:"index" json:"userId"`
	Action            string              `gorm:"not null" json:"action"` // "view_info", "get_transcript", "translate"
	SourceLanguage    string              `json:"sourceLanguage,omitempty"`
	TargetLanguage    string              `json:"targetLanguage,omitempty"`
	ProcessingTimeMs  int64               `json:"processingTimeMs,omitempty"`
	Success           bool                `gorm:"default:true" json:"success"`
	ErrorMessage      string              `json:"errorMessage,omitempty"`
	
	// Relationship
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName overrides the table name for VideoTranscriptCache
func (VideoTranscriptCache) TableName() string {
	return "video_transcript_cache"
}

// TableName overrides the table name for VideoTranslationCache
func (VideoTranslationCache) TableName() string {
	return "video_translation_cache"
}

// TableName overrides the table name for UserAPIKey
func (UserAPIKey) TableName() string {
	return "user_api_keys"
}

// TableName overrides the table name for VideoAnalytics
func (VideoAnalytics) TableName() string {
	return "video_analytics"
}