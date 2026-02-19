package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gino/cars-crud/pkg/config"
)

type RedisCache struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
	db, _ := strconv.Atoi(cfg.RedisDB)

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		Client: client,
		TTL:    5 * time.Minute,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(ctx context.Context, key string, value string) error {
	return c.Client.Set(ctx, key, value, c.TTL).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	iter := c.Client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		c.Client.Del(ctx, iter.Val())
	}
	return iter.Err()
}
