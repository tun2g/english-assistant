package docs

import (
	"app-backend/internal/dto"
	"app-backend/internal/types"
	"github.com/gin-gonic/gin"
)

// NewUserDocs creates instances of user-related DTOs for swagger documentation
// This function is never called but ensures the DTOs are considered "used" by the linter
func NewUserDocs() {
	_ = dto.UserResponse{}
	_ = dto.UpdateProfileRequest{}
	_ = dto.ChangePasswordRequest{}
	_ = dto.UserListResponse{}
	_ = types.PaginationMetadata{}
}

// UserGetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.UserResponse "User profile information"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [get]
func UserGetProfile(c *gin.Context) {}

// UserUpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.UpdateProfileRequest true "Profile update request"
// @Success 200 {object} dto.UserResponse "Updated user profile"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/profile [put]
func UserUpdateProfile(c *gin.Context) {}

// UserChangePassword godoc
// @Summary Change user password
// @Description Change the authenticated user's password
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Current password is incorrect"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/change-password [post]
func UserChangePassword(c *gin.Context) {}

// UserDeleteAccount godoc
// @Summary Delete user account
// @Description Delete the authenticated user's account (soft delete)
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]interface{} "Account deleted successfully"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/account [delete]
func UserDeleteAccount(c *gin.Context) {}

// UserListUsers godoc
// @Summary List users (Admin only)
// @Description Get a paginated list of users - requires admin role
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term for user name or email"
// @Success 200 {object} dto.UserListResponse "Paginated list of users"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /user/list [get]
func UserListUsers(c *gin.Context) {}