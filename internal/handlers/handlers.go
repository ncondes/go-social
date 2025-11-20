package handlers

import "github.com/ncondes/go/social/internal/services"

type Handlers struct {
	HealthHandler *HealthHandler
	UserHandler   *UserHandler
	PostHandler   *PostHandler
}

func New(
	services *services.Services,
) *Handlers {
	return &Handlers{
		HealthHandler: NewHealthHandler(),
		UserHandler:   NewUserHandler(services.UserService),
		PostHandler:   NewPostHandler(services.PostService),
	}
}
