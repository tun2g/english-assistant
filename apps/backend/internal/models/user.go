package models

type User struct {
	Auditable
	
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"-" gorm:"not null"` // Hidden from JSON responses
	Avatar    string `json:"avatar"`
	IsActive  bool   `json:"is_active" gorm:"default:true"`
	Role      string `json:"role" gorm:"default:'user'"`
}

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
	Avatar    *string `json:"avatar,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	Role      *string `json:"role,omitempty"`
}