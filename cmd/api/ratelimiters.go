package main

import (
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/ratelimit"
	"github.com/redis/go-redis/v9"
)

type RateLimiters struct {
	// Server protection
	Global ratelimit.RateLimiter

	// IP-based (anonymous users)
	StrictIP   ratelimit.RateLimiter
	ModerateIP ratelimit.RateLimiter

	// Operation-based (authenticated users)
	ReadOps  ratelimit.RateLimiter
	WriteOps ratelimit.RateLimiter
}

func newRateLimiters(rc *redis.Client, cfg config.RateLimitConfig) *RateLimiters {
	if !cfg.Enabled || rc == nil {
		return &RateLimiters{
			Global:     nil,
			StrictIP:   nil,
			ModerateIP: nil,
			ReadOps:    nil,
			WriteOps:   nil,
		}
	}

	return &RateLimiters{
		Global:     ratelimit.NewRedisRateLimiter(rc, cfg.Global),
		StrictIP:   ratelimit.NewRedisRateLimiter(rc, cfg.StrictIP),
		ModerateIP: ratelimit.NewRedisRateLimiter(rc, cfg.ModerateIP),
		ReadOps:    ratelimit.NewRedisRateLimiter(rc, cfg.ReadOps),
		WriteOps:   ratelimit.NewRedisRateLimiter(rc, cfg.WriteOps),
	}
}
