package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/mailer"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository     domain.UserRepositoryInterface
	followerRepository domain.FollowerRepositoryInterface
	roleRepository     domain.RoleRepositoryInterface
	config             *config.Config
	mailer             mailer.Mailer
	logger             logging.Logger
	authenticator      *auth.JWTAuthenticator
}

func NewUserService(
	userRepository domain.UserRepositoryInterface,
	followerRepository domain.FollowerRepositoryInterface,
	roleRepository domain.RoleRepositoryInterface,
	config *config.Config,
	mailer mailer.Mailer,
	logger logging.Logger,
	authenticator *auth.JWTAuthenticator,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		followerRepository: followerRepository,
		roleRepository:     roleRepository,
		config:             config,
		mailer:             mailer,
		logger:             logger,
		authenticator:      authenticator,
	}
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

func (s *UserService) RegisterUserWithInvitation(ctx context.Context, registerUserInput *domain.RegisterUserInput) (*domain.User, string, error) {
	defaultRole, err := s.roleRepository.GetByName(ctx, "user")
	if err != nil {
		return nil, "", err
	}

	user := domain.User{
		FirstName: registerUserInput.FirstName,
		LastName:  registerUserInput.LastName,
		Username:  registerUserInput.Username,
		Email:     registerUserInput.Email,
		RoleID:    defaultRole.ID,
	}
	// Hash the user's password
	if err := user.HashPassword(registerUserInput.Password); err != nil {
		return nil, "", err
	}
	// Generate a random token and hash it
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Execute transaction
	if err := s.userRepository.CreateUserAndInvitation(ctx, &user, registerUserInput.InvitationMethod, hashToken); err != nil {
		return nil, "", err
	}

	isSandbox := s.config.Env != "production"
	activationURL := fmt.Sprintf("%s/activate?token=%s", s.config.FrontendURL, plainToken)

	if err := s.mailer.Send(ctx, mailer.UserInvitationTemplate, mailer.Mail{
		To: []mailer.To{
			{
				Name:  user.FirstName,
				Email: user.Email,
			},
		},
		Args: mailer.InvitationEmailData{
			To: mailer.To{
				Name:  user.FirstName,
				Email: user.Email,
			},
			ActivationURL: activationURL,
		},
	}, isSandbox); err != nil {
		s.logger.Errorw("Failed to send invitation email", "error", err)
		// Rollback user creation if email sending fails (SAGA pattern)
		if err := s.userRepository.DeleteUser(ctx, user.ID); err != nil {
			s.logger.Errorw("Failed to delete user after email sending failed", "error", err)
		}

		return nil, "", err
	}

	return &user, plainToken, nil
}

func (s *UserService) ActivateUser(ctx context.Context, token string) error {
	return s.userRepository.ActivateUser(ctx, token)
}

func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if err := s.verifyPassword(user, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.authenticator.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) verifyPassword(user *domain.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return err
	}

	return nil
}
