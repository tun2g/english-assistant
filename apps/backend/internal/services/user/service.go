package user

import (
	"app-backend/internal/dto"
	"app-backend/internal/errors"
	"app-backend/internal/models"
	"app-backend/internal/repositories"
	"app-backend/internal/types"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	userRepo repositories.UserRepositoryInterface
}

func NewUserService(userRepo repositories.UserRepositoryInterface) ServiceInterface {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) CreateUser(req *dto.RegisterRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.NewAppError("Failed to check existing user", err, http.StatusInternalServerError)
	}
	if existingUser != nil {
		return nil, errors.NewAppError("User already exists", nil, http.StatusConflict)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewAppError("Failed to hash password", err, http.StatusInternalServerError)
	}

	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		Role:      "user", // Default role
		IsActive:  true,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, errors.NewAppError("Failed to create user", err, http.StatusInternalServerError)
	}

	return user, nil
}

func (s *Service) GetUser(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("User not found", err, http.StatusNotFound)
		}
		return nil, errors.NewAppError("Failed to get user", err, http.StatusInternalServerError)
	}
	return user, nil
}

func (s *Service) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("User not found", err, http.StatusNotFound)
		}
		return nil, errors.NewAppError("Failed to get user", err, http.StatusInternalServerError)
	}
	return user, nil
}

func (s *Service) UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("User not found", err, http.StatusNotFound)
		}
		return nil, errors.NewAppError("Failed to get user", err, http.StatusInternalServerError)
	}

	// Update fields if provided
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Email != nil {
		// Check if email already exists
		existingUser, err := s.userRepo.GetByEmail(*req.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, errors.NewAppError("Failed to check existing email", err, http.StatusInternalServerError)
		}
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.NewAppError("Email already in use", nil, http.StatusConflict)
		}
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, errors.NewAppError("Failed to update user", err, http.StatusInternalServerError)
	}

	return user, nil
}

func (s *Service) DeleteUser(id uint) error {
	_, err := s.userRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError("User not found", err, http.StatusNotFound)
		}
		return errors.NewAppError("Failed to get user", err, http.StatusInternalServerError)
	}

	err = s.userRepo.Delete(id)
	if err != nil {
		return errors.NewAppError("Failed to delete user", err, http.StatusInternalServerError)
	}

	return nil
}

func (s *Service) ListUsers(pagReq *types.PaginationRequest) (*types.PaginationResponse[models.User], error) {
	users, err := s.userRepo.List(pagReq, nil)
	if err != nil {
		return nil, errors.NewAppError("Failed to list users", err, http.StatusInternalServerError)
	}
	return users, nil
}

func (s *Service) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewAppError("User not found", err, http.StatusNotFound)
		}
		return errors.NewAppError("Failed to get user", err, http.StatusInternalServerError)
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
	if err != nil {
		return errors.NewAppError("Invalid current password", nil, http.StatusBadRequest)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewAppError("Failed to hash password", err, http.StatusInternalServerError)
	}

	user.Password = string(hashedPassword)
	err = s.userRepo.Update(user)
	if err != nil {
		return errors.NewAppError("Failed to update password", err, http.StatusInternalServerError)
	}

	return nil
}