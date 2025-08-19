package models

import (
	"time"
)

// Session represents a user session in the database
type Session struct {
	Auditable
	
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	TokenHash string    `json:"-" gorm:"uniqueIndex;not null"` // JWT token hash for validation
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	LastUsed  time.Time `json:"last_used"`
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is active and not expired
func (s *Session) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}