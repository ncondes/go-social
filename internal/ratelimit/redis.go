package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/ncondes/go/social/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client *redis.Client
	config config.RateLimitTier
}

func NewRedisRateLimiter(client *redis.Client, config config.RateLimitTier) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		config: config,
	}
}

func (rl *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, *LimitInfo, error) {
	// Redis key with time window
	windowKey := fmt.Sprintf("ratelimit:%s:%d", key, time.Now().Unix()/int64(rl.config.Window.Seconds()))
	// Increment counter
	count, err := rl.client.Incr(ctx, windowKey).Result()
	if err != nil {
		return false, nil, err
	}
	// Set expiration on first request
	if count == 1 {
		rl.client.Expire(ctx, windowKey, rl.config.Window)
	}
	// Calculate metadata
	remaining := rl.config.RequestsPerWindow - int(count)
	if remaining < 0 {
		remaining = 0
	}

	reset := time.Now().Add(rl.config.Window)
	allowed := count <= int64(rl.config.RequestsPerWindow)

	info := &LimitInfo{
		Limit:      rl.config.RequestsPerWindow,
		Remaining:  remaining,
		Reset:      reset,
		RetryAfter: time.Duration(rl.config.Window.Seconds()),
	}

	return allowed, info, nil
}
