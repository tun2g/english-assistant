package user

import (
	"app-backend/internal/dto"
	"app-backend/internal/models"
	"app-backend/internal/types"
)

type ServiceInterface interface {
	CreateUser(req *dto.RegisterRequest) (*models.User, error)
	GetUser(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers(pagReq *types.PaginationRequest) (*types.PaginationResponse[models.User], error)
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
}