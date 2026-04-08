package cache

import (
	"time"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	queryTimeoutDuration = 5 * time.Second
)

type Storage struct {
	UserStorage domain.UserStorageInterface
}

func NewCacheStorage(client *redis.Client) *Storage {
	if client == nil {
		return nil
	}

	return &Storage{
		UserStorage: NewUserStorage(client),
	}
}
