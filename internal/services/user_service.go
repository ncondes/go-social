package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type UserService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepository.Create(ctx, user)
}
