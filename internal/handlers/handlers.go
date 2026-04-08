package handlers

import (
	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/services"
)

type Handlers struct {
	HealthHandler  *HealthHandler
	UserHandler    *UserHandler
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
	FeedHandler    *FeedHandler
	AuthHandler    *AuthHandler
}

func New(
	config *config.Config,
	services *services.Services,
	validator *Validator,
	logger logging.Logger,
	authorizer *auth.Authorizer,
) *Handlers {
	return &Handlers{
		HealthHandler:  NewHealthHandler(config, logger),
		UserHandler:    NewUserHandler(services.UserService, logger),
		PostHandler:    NewPostHandler(services.PostService, validator, logger, authorizer),
		CommentHandler: NewCommentHandler(services.CommentService, validator, logger),
		FeedHandler:    NewFeedHandler(services.FeedService, logger),
		AuthHandler:    NewAuthHandler(services.UserService, validator, logger),
	}
}
