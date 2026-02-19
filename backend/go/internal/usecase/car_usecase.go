package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/gino/cars-crud/internal/cache"
	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/repository"
)

type CarUsecase struct {
	repo  repository.CarRepository
	cache *cache.RedisCache
}

func NewCarUsecase(repo repository.CarRepository, cache *cache.RedisCache) *CarUsecase {
	return &CarUsecase{repo: repo, cache: cache}
}

func (u *CarUsecase) Create(ctx context.Context, req domain.CreateCarRequest) (*domain.Car, error) {
	car := &domain.Car{
		Brand: req.Brand,
		Model: req.Model,
		Year:  req.Year,
		Color: req.Color,
		Price: req.Price,
	}

	if err := u.repo.Create(ctx, car); err != nil {
		return nil, err
	}

	_ = u.cache.DeleteByPattern(ctx, "cars:list:*")
	return car, nil
}

func (u *CarUsecase) GetByID(ctx context.Context, id uuid.UUID) (*domain.Car, error) {
	key := fmt.Sprintf("cars:%s", id.String())

	cached, err := u.cache.Get(ctx, key)
	if err == nil {
		var car domain.Car
		if json.Unmarshal([]byte(cached), &car) == nil {
			return &car, nil
		}
	} else if err != redis.Nil {
		return nil, err
	}

	car, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(car); err == nil {
		_ = u.cache.Set(ctx, key, string(data))
	}

	return car, nil
}

func (u *CarUsecase) GetAll(ctx context.Context, offset, limit int) ([]domain.Car, int64, error) {
	key := fmt.Sprintf("cars:list:%d:%d", offset, limit)

	type listCache struct {
		Cars  []domain.Car `json:"cars"`
		Total int64        `json:"total"`
	}

	cached, err := u.cache.Get(ctx, key)
	if err == nil {
		var lc listCache
		if json.Unmarshal([]byte(cached), &lc) == nil {
			return lc.Cars, lc.Total, nil
		}
	}

	cars, total, err := u.repo.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	if data, err := json.Marshal(listCache{Cars: cars, Total: total}); err == nil {
		_ = u.cache.Set(ctx, key, string(data))
	}

	return cars, total, nil
}

func (u *CarUsecase) Update(ctx context.Context, id uuid.UUID, req domain.UpdateCarRequest) (*domain.Car, error) {
	car, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Brand != nil {
		car.Brand = *req.Brand
	}
	if req.Model != nil {
		car.Model = *req.Model
	}
	if req.Year != nil {
		car.Year = *req.Year
	}
	if req.Color != nil {
		car.Color = *req.Color
	}
	if req.Price != nil {
		car.Price = *req.Price
	}

	if err := u.repo.Update(ctx, car); err != nil {
		return nil, err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("cars:%s", id.String()))
	_ = u.cache.DeleteByPattern(ctx, "cars:list:*")

	return car, nil
}

func (u *CarUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("cars:%s", id.String()))
	_ = u.cache.DeleteByPattern(ctx, "cars:list:*")

	return nil
}
