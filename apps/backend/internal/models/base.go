package models

import (
	"time"

	"gorm.io/gorm"
)

// Auditable contains common fields for all models
type Auditable struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to set creation timestamp
func (a *Auditable) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate hook to set update timestamp
func (a *Auditable) BeforeUpdate(tx *gorm.DB) error {
	a.UpdatedAt = time.Now()
	return nil
}