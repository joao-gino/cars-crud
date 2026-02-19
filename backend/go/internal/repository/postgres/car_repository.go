package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/repository"
)

type carRepository struct {
	db *gorm.DB
}

func NewCarRepository(db *gorm.DB) repository.CarRepository {
	return &carRepository{db: db}
}

func (r *carRepository) Create(ctx context.Context, car *domain.Car) error {
	return r.db.WithContext(ctx).Create(car).Error
}

func (r *carRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	var car domain.Car
	if err := r.db.WithContext(ctx).First(&car, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

func (r *carRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.Car, int64, error) {
	var cars []domain.Car
	var total int64

	if err := r.db.WithContext(ctx).Model(&domain.Car{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("created_at DESC").Find(&cars).Error; err != nil {
		return nil, 0, err
	}

	return cars, total, nil
}

func (r *carRepository) Update(ctx context.Context, car *domain.Car) error {
	return r.db.WithContext(ctx).Save(car).Error
}

func (r *carRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Car{}, "id = ?", id).Error
}
