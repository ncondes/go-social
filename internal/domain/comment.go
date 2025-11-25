package domain

import (
	"context"
	"errors"
	"time"

	"github.com/ncondes/go/social/internal/dtos"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetManyByPostID(ctx context.Context, postID int64) ([]*dtos.CommentResponseDTO, error)
	GetCountByPostID(ctx context.Context, postID int64) (int, error)
}

var (
	ErrCommentNotFound = errors.New("comment not found")
)
