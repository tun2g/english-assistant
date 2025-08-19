package user

import "github.com/gin-gonic/gin"

// HandlerInterface defines the contract for user handlers
type HandlerInterface interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	ChangePassword(c *gin.Context)
	DeleteAccount(c *gin.Context)
	ListUsers(c *gin.Context)
}