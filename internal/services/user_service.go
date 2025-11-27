package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type UserService struct {
	userRepository     domain.UserRepository
	followerRepository domain.FollowerRepository
}

func NewUserService(userRepository domain.UserRepository, followerRepository domain.FollowerRepository) *UserService {
	return &UserService{userRepository: userRepository, followerRepository: followerRepository}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.Create(ctx, user)
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
