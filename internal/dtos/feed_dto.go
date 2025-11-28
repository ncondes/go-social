package dtos

import (
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

type FeedPostResponseDTO struct {
	ID           int64         `json:"id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Tags         []string      `json:"tags"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Author       AuthorInfoDTO `json:"author"`
	CommentCount int           `json:"comment_count"`
}

func (dto *FeedPostResponseDTO) FromDomain(feedPost *domain.FeedPost) *FeedPostResponseDTO {
	response := &FeedPostResponseDTO{
		ID:        feedPost.Post.ID,
		Title:     feedPost.Post.Title,
		Content:   feedPost.Post.Content,
		Tags:      feedPost.Post.Tags,
		CreatedAt: feedPost.Post.CreatedAt,
		UpdatedAt: feedPost.Post.UpdatedAt,
		Author: AuthorInfoDTO{
			ID:       feedPost.Author.ID,
			Username: feedPost.Author.Username,
			Fullname: feedPost.Author.FullName(),
		},
		CommentCount: feedPost.CommentCount,
	}

	return response
}

type FeedResponseDTO struct {
	Posts      []*FeedPostResponseDTO    `json:"posts"`
	Pagination CursorBasedPaginationMeta `json:"pagination"`
}
