package dtos

import "time"

type CreatePostDTO struct {
	Title   string   `json:"title"   validate:"required,min=1,max=100"`
	Content string   `json:"content" validate:"required,min=1,max=1000"`
	Tags    []string `json:"tags"    validate:"omitempty,dive,min=1,max=50"`
}

type UpdatePostDTO struct {
	Title     *string    `json:"title"      validate:"omitempty,min=1,max=100"`
	Content   *string    `json:"content"    validate:"omitempty,min=1,max=1000"`
	Tags      *[]string  `json:"tags"       validate:"omitempty,dive,min=1,max=50"`
	UpdatedAt *time.Time `json:"updated_at" validate:"required"`
}

type PostResponseDTO struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Tags         []string      `json:"tags"`
	CommentCount int           `json:"comment_count"`
	Author       AuthorInfoDTO `json:"author"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}
