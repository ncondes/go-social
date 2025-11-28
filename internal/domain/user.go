package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // When marshaling, don't include this field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (*User, error)
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (*User, error)
	FollowUser(ctx context.Context, userID int64, followerID int64) error
	UnfollowUser(ctx context.Context, userID int64, followerID int64) error
}

var (
	ErrUserNotFound = errors.New("user not found")
)
