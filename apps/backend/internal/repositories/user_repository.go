package repositories

import (
	"app-backend/internal/models"
	"app-backend/internal/types"
	
	"gorm.io/gorm"
)

// UserRepositoryInterface extends base repository with user-specific methods
type UserRepositoryInterface interface {
	BaseRepositoryInterface[models.User]
	GetByEmail(email string) (*models.User, error)
	GetActiveUsers(req *types.PaginationRequest) (*types.PaginationResponse[models.User], error)
}

// UserRepository implements user-specific repository
type UserRepository struct {
	*BaseRepository[models.User]
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{
		BaseRepository: NewBaseRepository[models.User](db),
	}
}

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	return r.FindBy("email", email)
}

// GetActiveUsers retrieves only active users with pagination
func (r *UserRepository) GetActiveUsers(req *types.PaginationRequest) (*types.PaginationResponse[models.User], error) {
	opts := &QueryOptions{
		Conditions: map[string]interface{}{
			"is_active": true,
		},
		SearchFields: []string{"first_name", "last_name", "email"},
	}
	return r.List(req, opts)
}