package services

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
)

type CommentService struct {
	commentRepository domain.CommentRepository
}

func NewCommentService(commentRepository domain.CommentRepository) *CommentService {
	return &CommentService{commentRepository: commentRepository}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	return s.commentRepository.Create(ctx, comment)
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64) ([]*dtos.CommentResponseDTO, error) {
	return s.commentRepository.GetManyByPostID(ctx, postID)
}
