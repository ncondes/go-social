package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type CommentService struct {
	commentRepository domain.CommentRepositoryInterface
}

func NewCommentService(commentRepository domain.CommentRepositoryInterface) *CommentService {
	return &CommentService{commentRepository: commentRepository}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	return s.commentRepository.Create(ctx, comment)
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64) ([]*domain.CommentWithAuthor, error) {
	return s.commentRepository.GetManyByPostID(ctx, postID)
}
