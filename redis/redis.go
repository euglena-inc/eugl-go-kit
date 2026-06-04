package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func New(ctx context.Context, cfg Config) (*goredis.Client, error) {
	if cfg.Addr == "" {
		return nil, nil
	}

	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}

func SetWithTTL(ctx context.Context, client *goredis.Client, key string, value interface{}, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("redis ttl must be positive")
	}
	return client.Set(ctx, key, value, ttl).Err()
}

func SetNXWithTTL(ctx context.Context, client *goredis.Client, key string, value interface{}, ttl time.Duration) (bool, error) {
	if ttl <= 0 {
		return false, fmt.Errorf("redis ttl must be positive")
	}
	return client.SetNX(ctx, key, value, ttl).Result()
}
