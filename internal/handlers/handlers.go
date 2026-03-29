package handlers

import (
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/services"
)

type Handlers struct {
	HealthHandler  *HealthHandler
	UserHandler    *UserHandler
	PostHandler    *PostHandler
	CommentHandler *CommentHandler
	FeedHandler    *FeedHandler
}

func New(
	config *config.Config,
	services *services.Services,
	validator *Validator,
) *Handlers {
	return &Handlers{
		HealthHandler:  NewHealthHandler(config),
		UserHandler:    NewUserHandler(services.UserService),
		PostHandler:    NewPostHandler(services.PostService, validator),
		CommentHandler: NewCommentHandler(services.CommentService, validator),
		FeedHandler:    NewFeedHandler(services.FeedService),
	}
}
