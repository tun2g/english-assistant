package repositories

import (
	"app-backend/internal/models"
	"app-backend/internal/types"
	"time"

	"gorm.io/gorm"
)

type SessionRepositoryInterface interface {
	BaseRepositoryInterface[models.Session]
	GetByTokenHash(tokenHash string) (*models.Session, error)
	GetActiveSessionsByUserID(userID uint) ([]*models.Session, error)
	DeactivateSession(sessionID uint) error
	DeactivateUserSessions(userID uint) error
	CleanupExpiredSessions() error
	UpdateLastUsed(sessionID uint) error
}

type SessionRepository struct {
	*BaseRepository[models.Session]
}

func NewSessionRepository(db *gorm.DB) SessionRepositoryInterface {
	return &SessionRepository{
		BaseRepository: NewBaseRepository[models.Session](db),
	}
}

// GetByTokenHash finds a session by its token hash
func (r *SessionRepository) GetByTokenHash(tokenHash string) (*models.Session, error) {
	var session models.Session
	err := r.GetDB().Preload("User").Where("token_hash = ? AND is_active = ?", tokenHash, true).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetActiveSessionsByUserID retrieves all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(userID uint) ([]*models.Session, error) {
	opts := &QueryOptions{
		Conditions: map[string]interface{}{
			"user_id":   userID,
			"is_active": true,
		},
	}
	
	req := &types.PaginationRequest{
		Page:     1,
		PageSize: 100, // Get all active sessions
		SortBy:   "last_used",
		SortDir:  "desc",
	}
	
	response, err := r.List(req, opts)
	if err != nil {
		return nil, err
	}
	
	// Convert to slice of pointers
	sessions := make([]*models.Session, len(response.Data))
	for i := range response.Data {
		sessions[i] = &response.Data[i]
	}
	
	return sessions, nil
}

// DeactivateSession marks a session as inactive
func (r *SessionRepository) DeactivateSession(sessionID uint) error {
	return r.GetDB().Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("is_active", false).Error
}

// DeactivateUserSessions marks all user sessions as inactive
func (r *SessionRepository) DeactivateUserSessions(userID uint) error {
	return r.GetDB().Model(&models.Session{}).
		Where("user_id = ?", userID).
		Update("is_active", false).Error
}

// CleanupExpiredSessions removes expired sessions from database
func (r *SessionRepository) CleanupExpiredSessions() error {
	return r.GetDB().Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

// UpdateLastUsed updates the last used timestamp for a session
func (r *SessionRepository) UpdateLastUsed(sessionID uint) error {
	return r.GetDB().Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("last_used", time.Now()).Error
}