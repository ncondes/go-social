package domain

import (
	"context"
	"errors"
	"time"
)

type Follower struct {
	UserID     int64     `json:"user_id"`
	FollowerID int64     `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerRepositoryInterface interface {
	FollowUser(ctx context.Context, userID int64, followerID int64) error
	UnfollowUser(ctx context.Context, userID int64, followerID int64) error
}

var (
	ErrUserAlreadyFollowing = errors.New("user already following")
)
