package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type PostService struct {
	postRepository domain.PostRepositoryInterface
}

func NewPostService(postRepository domain.PostRepositoryInterface) *PostService {
	return &PostService{postRepository: postRepository}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	return s.postRepository.Create(ctx, post)
}

func (s *PostService) GetPost(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
	return s.postRepository.GetByID(ctx, postID)
}

func (s *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	// TODO: should we validate the post before updating?
	// TODO: maybe caching the post and getting the values from it
	return s.postRepository.Update(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, postID int64) error {
	return s.postRepository.Delete(ctx, postID)
}
