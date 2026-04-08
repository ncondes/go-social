package domain

import (
	"context"
	"errors"
	"time"
)

type Post struct {
	ID        int64     `json:"id"         example:"1"`
	Title     string    `json:"title"      example:"My first post"`
	Content   string    `json:"content"    example:"This is the content of my post."`
	UserID    int64     `json:"user_id"    example:"1"`
	Tags      []string  `json:"tags"       example:"go,api"`
	Version   int64     `json:"version"    example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-03T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-04-03T10:00:00Z"`
}

func (p *Post) GetOwnerID() any {
	return p.UserID
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
	ErrPostNotFound        = errors.New("post not found")
	ErrPostVersionConflict = errors.New("post version conflict")
)
