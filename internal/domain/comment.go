package domain

import (
	"context"
	"errors"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentWithAuthor struct {
	Comment
	Author User
}

type CommentRepositoryInterface interface {
	Create(ctx context.Context, comment *Comment) error
	GetManyByPostID(ctx context.Context, postID int64) ([]*CommentWithAuthor, error)
	GetCountByPostID(ctx context.Context, postID int64) (int, error)
}

type CommentServiceInterface interface {
	CreateComment(ctx context.Context, comment *Comment) error
	GetCommentsByPostID(ctx context.Context, postID int64) ([]*CommentWithAuthor, error)
}

var (
	ErrCommentNotFound = errors.New("comment not found")
)
