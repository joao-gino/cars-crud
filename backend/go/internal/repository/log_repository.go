package repository

import (
	"context"

	"github.com/gino/cars-crud/internal/domain"
)

type LogRepository interface {
	GetAll(ctx context.Context, offset, limit int) ([]domain.RequestLog, int64, error)
}
