package domain

import (
	"context"
	"errors"
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostWithDetails struct {
	Post
	Author       User
	CommentCount int
}

type PostRepositoryInterface interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, postID int64) (*PostWithDetails, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, postID int64) error
}

type PostServiceInterface interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPost(ctx context.Context, postID int64) (*PostWithDetails, error)
	UpdatePost(ctx context.Context, post *Post) error
	DeletePost(ctx context.Context, postID int64) error
}

var (
	ErrPostNotFound = errors.New("post not found")
)
