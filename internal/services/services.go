package services

import (
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/mailer"
	"github.com/ncondes/go/social/internal/repositories"
)

type Services struct {
	UserService    *UserService
	PostService    *PostService
	CommentService *CommentService
	FeedService    *FeedService
}

func New(
	repositories *repositories.Repositories,
	config *config.Config,
	mailer mailer.Mailer,
	logger logging.Logger,
) *Services {
	return &Services{
		UserService: NewUserService(
			repositories.UserRepository,
			repositories.FollowerRepository,
			config,
			mailer,
			logger,
		),
		PostService:    NewPostService(repositories.PostRepository),
		CommentService: NewCommentService(repositories.CommentRepository),
		FeedService:    NewFeedService(repositories.FeedRepository),
	}
}
