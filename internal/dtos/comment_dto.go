package dtos

import "time"

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
