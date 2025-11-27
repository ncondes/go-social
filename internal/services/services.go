package services

import (
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
) *Services {
	return &Services{
		UserService:    NewUserService(repositories.UserRepository, repositories.FollowerRepository),
		PostService:    NewPostService(repositories.PostRepository),
		CommentService: NewCommentService(repositories.CommentRepository),
		FeedService:    NewFeedService(repositories.FeedRepository),
	}
}
