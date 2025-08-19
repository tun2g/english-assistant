package auth

import "github.com/gin-gonic/gin"

// HandlerInterface defines the contract for authentication handlers
type HandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	LogoutAll(c *gin.Context)
	RefreshToken(c *gin.Context)
	GetSessions(c *gin.Context)
	RevokeSession(c *gin.Context)
}