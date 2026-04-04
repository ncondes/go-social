package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/ncondes/go/social/internal/domain"
)

type UserService struct {
	userRepository     domain.UserRepositoryInterface
	followerRepository domain.FollowerRepositoryInterface
}

func NewUserService(userRepository domain.UserRepositoryInterface, followerRepository domain.FollowerRepositoryInterface) *UserService {
	return &UserService{userRepository: userRepository, followerRepository: followerRepository}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.CreateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	return s.userRepository.GetUser(ctx, id)
}

func (s *UserService) FollowUser(ctx context.Context, userID int64, followerID int64) error {
	return s.followerRepository.FollowUser(ctx, userID, followerID)
}

func (s *UserService) UnfollowUser(ctx context.Context, userID int64, followerID int64) error {
	return s.followerRepository.UnfollowUser(ctx, userID, followerID)
}

func (s *UserService) RegisterUserWithInvitation(ctx context.Context, registerUserInput *domain.RegisterUserInput) (*domain.UserWithInvitationToken, error) {
	user := domain.User{
		FirstName: registerUserInput.FirstName,
		LastName:  registerUserInput.LastName,
		Username:  registerUserInput.Username,
		Email:     registerUserInput.Email,
	}
	// Hash the user's password
	if err := user.HashPassword(registerUserInput.Password); err != nil {
		return nil, err
	}
	// Generate a random token and hash it
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Execute transaction
	if err := s.userRepository.CreateUserAndInvitation(ctx, &user, registerUserInput.InvitationMethod, hashToken); err != nil {
		return nil, err
	}

	// TODO: Send the invitation (email or SMS)

	return &domain.UserWithInvitationToken{
		User:  user,
		Token: plainToken,
	}, nil
}

func (s *UserService) ActivateUser(ctx context.Context, token string) error {
	return s.userRepository.ActivateUser(ctx, token)
}
