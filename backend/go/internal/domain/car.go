package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Car struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey" example:"550e8400-e29b-41d4-a716-446655440000"`
	Brand     string         `json:"brand" gorm:"not null;size:100" example:"Toyota"`
	Model     string         `json:"model" gorm:"not null;size:100" example:"Corolla"`
	Year      int            `json:"year" gorm:"not null" example:"2024"`
	Color     string         `json:"color" gorm:"not null;size:50" example:"White"`
	Price     float64        `json:"price" gorm:"not null" example:"35000.00"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateCarRequest struct {
	Brand string  `json:"brand" example:"Toyota"`
	Model string  `json:"model" example:"Corolla"`
	Year  int     `json:"year" example:"2024"`
	Color string  `json:"color" example:"White"`
	Price float64 `json:"price" example:"35000.00"`
}

type UpdateCarRequest struct {
	Brand *string  `json:"brand,omitempty" example:"Honda"`
	Model *string  `json:"model,omitempty" example:"Civic"`
	Year  *int     `json:"year,omitempty" example:"2025"`
	Color *string  `json:"color,omitempty" example:"Black"`
	Price *float64 `json:"price,omitempty" example:"40000.00"`
}

func (c *Car) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
