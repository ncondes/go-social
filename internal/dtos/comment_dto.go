package dtos

import (
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

type CreateCommentDTO struct {
	Content string `json:"content" validate:"required,max=1000"`
}

type CommentResponseDTO struct {
	ID        int64         `json:"id"`
	PostID    int64         `json:"post_id"`
	Author    AuthorInfoDTO `json:"author"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func (dto *CommentResponseDTO) FromDomain(comment *domain.CommentWithAuthor) *CommentResponseDTO {
	return &CommentResponseDTO{
		ID:      comment.ID,
		PostID:  comment.PostID,
		Content: comment.Content,
		Author: AuthorInfoDTO{
			ID:       comment.Author.ID,
			Username: comment.Author.Username,
			Fullname: comment.Author.FullName(),
		},
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
