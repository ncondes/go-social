package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type PostService struct {
	postRepository domain.PostRepository
}

func NewPostService(postRepository domain.PostRepository) *PostService {
	return &PostService{postRepository: postRepository}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	return s.postRepository.Create(ctx, post)
}
