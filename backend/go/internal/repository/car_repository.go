package repository

import (
	"context"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/google/uuid"
)

type CarRepository interface {
	Create(ctx context.Context, car *domain.Car) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error)
	GetAll(ctx context.Context, offset, limit int) ([]domain.Car, int64, error)
	Update(ctx context.Context, car *domain.Car) error
	Delete(ctx context.Context, id uuid.UUID) error
}
