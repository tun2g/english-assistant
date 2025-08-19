package types

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// UserContext represents the authenticated user context
type UserContext struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	SessionID uint   `json:"session_id"`
}

// GetUserContext extracts user context from gin.Context
func GetUserContext(c *gin.Context) (*UserContext, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("user_id not found in context")
	}

	email, exists := c.Get("user_email")
	if !exists {
		return nil, errors.New("user_email not found in context")
	}

	role, exists := c.Get("user_role")
	if !exists {
		return nil, errors.New("user_role not found in context")
	}

	sessionID, exists := c.Get("session_id")
	if !exists {
		return nil, errors.New("session_id not found in context")
	}

	return &UserContext{
		UserID:    userID.(uint),
		Email:     email.(string),
		Role:      role.(string),
		SessionID: sessionID.(uint),
	}, nil
}

// SetUserContext sets user context in gin.Context
func SetUserContext(c *gin.Context, userCtx *UserContext) {
	c.Set("user_id", userCtx.UserID)
	c.Set("user_email", userCtx.Email)
	c.Set("user_role", userCtx.Role)
	c.Set("session_id", userCtx.SessionID)
}

// IsAdmin checks if the user has admin role
func (uc *UserContext) IsAdmin() bool {
	return uc.Role == "admin"
}

// IsModerator checks if the user has moderator role
func (uc *UserContext) IsModerator() bool {
	return uc.Role == "moderator"
}

// HasRole checks if the user has any of the specified roles
func (uc *UserContext) HasRole(roles ...string) bool {
	for _, role := range roles {
		if uc.Role == role {
			return true
		}
	}
	return false
}

// HasPermission checks if the user has admin or moderator permissions
func (uc *UserContext) HasPermission() bool {
	return uc.IsAdmin() || uc.IsModerator()
}