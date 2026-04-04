package domain

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	Password  string    `json:"-"` // When marshaling, don't include this field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

type UserUpdate struct {
	FirstName *string
	LastName  *string
	Username  *string
	Email     *string
	IsActive  *bool
}

type UserWithInvitationToken struct {
	User
	Token string `json:"token"`
}

type RegisterUserInput struct {
	FirstName        string `json:"first_name" validate:"required,min=1,max=255" example:"John"`
	LastName         string `json:"last_name" validate:"required,min=1,max=255" example:"Doe"`
	Username         string `json:"username" validate:"required,min=1,max=255" example:"johndoe"`
	Email            string `json:"email" validate:"required,email" example:"johndoe@example.com"`
	Password         string `json:"password" validate:"required,min=8,max=255" example:"password123"`
	InvitationMethod string `json:"invitation_method" validate:"required,oneof=email sms" example:"email"`
}

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUserAndInvitation(ctx context.Context, user *User, method string, token string) error
	ActivateUser(ctx context.Context, token string) error
}

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id int64) (*User, error)
	FollowUser(ctx context.Context, userID int64, followerID int64) error
	UnfollowUser(ctx context.Context, userID int64, followerID int64) error
	RegisterUserWithInvitation(ctx context.Context, registerUserInput *RegisterUserInput) (*UserWithInvitationToken, error)
	ActivateUser(ctx context.Context, token string) error
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserEmailTaken     = errors.New("user email is already in use")
	ErrUserUsernameTaken  = errors.New("user username is already in use")
	ErrNoUserUpdateFields = errors.New("no fields to update on user")
)
