package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	UserExpirationTime = 60 * time.Minute
)

type UserStorage struct {
	client *redis.Client
}

func NewUserStorage(client *redis.Client) *UserStorage {
	if client == nil {
		return nil
	}

	return &UserStorage{client: client}
}

func (s *UserStorage) Get(ctx context.Context, userID int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	key := s.userKey(userID)

	data, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // cache miss - not an error
	}

	if err != nil {
		return nil, err
	}

	var user domain.User

	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) Set(ctx context.Context, userID int64, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	if user == nil {
		return domain.ErrNilUser
	}

	key := s.userKey(userID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, key, data, UserExpirationTime).Err()
}

func (s *UserStorage) Delete(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	key := s.userKey(userID)
	return s.client.Del(ctx, key).Err()
}

func (s *UserStorage) userKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}
