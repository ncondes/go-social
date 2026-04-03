package dtos

import (
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

type CreatePostDTO struct {
	Title   string   `json:"title"   validate:"required,min=1,max=100"       example:"My first post"`
	Content string   `json:"content" validate:"required,min=1,max=1000"      example:"This is the content of my post."`
	Tags    []string `json:"tags"    validate:"omitempty,dive,min=1,max=50"  example:"go,api"`
}

type UpdatePostDTO struct {
	Title   *string   `json:"title"   validate:"omitempty,min=1,max=100"    example:"My updated post"`
	Content *string   `json:"content" validate:"omitempty,min=1,max=1000"   example:"This is the updated content."`
	Tags    *[]string `json:"tags"    validate:"omitempty,dive,min=1,max=50" example:"go,api"`
	Version *int64    `json:"version" validate:"required"                   example:"1"`
}

type PostResponseDTO struct {
	ID           int64         `json:"id"            example:"1"`
	Title        string        `json:"title"         example:"My first post"`
	Content      string        `json:"content"       example:"This is the content of my post."`
	Tags         []string      `json:"tags"          example:"go,api"`
	Version      int64         `json:"version"       example:"1"`
	CommentCount int           `json:"comment_count" example:"3"`
	Author       AuthorInfoDTO `json:"author"`
	CreatedAt    time.Time     `json:"created_at"    example:"2026-04-03T10:00:00Z"`
	UpdatedAt    time.Time     `json:"updated_at"    example:"2026-04-03T10:00:00Z"`
}

func (dto *PostResponseDTO) FromDomain(post *domain.PostWithDetails) *PostResponseDTO {
	return &PostResponseDTO{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		Tags:         post.Tags,
		Version:      post.Version,
		CommentCount: post.CommentCount,
		Author: AuthorInfoDTO{
			ID:       post.Author.ID,
			Username: post.Author.Username,
			Fullname: post.Author.FullName(),
		},
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}
