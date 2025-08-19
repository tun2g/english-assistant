package user

import (
	"app-backend/internal/dto"
	"app-backend/internal/errors"
	"app-backend/internal/logger"
	"app-backend/internal/models"
	"app-backend/internal/services/user"
	"app-backend/internal/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	userService user.ServiceInterface
	logger      *logger.Logger
}

func NewUserHandler(userService user.ServiceInterface, logger *logger.Logger) HandlerInterface {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *Handler) GetProfile(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUser(userCtx.UserID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Get profile failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected get profile error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userResponse := &dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, userResponse)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid update profile request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	user, err := h.userService.UpdateUser(userCtx.UserID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Update profile failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected update profile error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userResponse := &dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	h.logger.Info("Profile updated successfully", zap.Uint("user_id", userCtx.UserID))
	c.JSON(http.StatusOK, userResponse)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid change password request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	err = h.userService.ChangePassword(userCtx.UserID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Change password failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected change password error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Password changed successfully", zap.Uint("user_id", userCtx.UserID))
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *Handler) DeleteAccount(c *gin.Context) {
	userCtx, err := types.GetUserContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.userService.DeleteUser(userCtx.UserID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("Delete account failed", zap.Error(err), zap.Uint("user_id", userCtx.UserID))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected delete account error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Account deleted successfully", zap.Uint("user_id", userCtx.UserID))
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *Handler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")
	search := c.Query("search")

	pagReq := &types.PaginationRequest{
		Page:     page,
		PageSize: pageSize,
		SortBy:   sortBy,
		SortDir:  sortDir,
		Search:   search,
	}

	users, err := h.userService.ListUsers(pagReq)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			h.logger.Error("List users failed", zap.Error(err))
			c.JSON(appErr.Status, gin.H{"error": appErr.Message})
			return
		}
		h.logger.Error("Unexpected list users error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, users)
}