package dtos

import "time"

type FeedPostResponseDTO struct {
	ID              int64         `json:"id"`
	Title           string        `json:"title"`
	Content         string        `json:"content"`
	Tags            []string      `json:"tags"`
	Author          AuthorInfoDTO `json:"author"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	CommentCount    int           `json:"comment_count"`
	RecencyScore    float64       `json:"-"`
	EngagementScore float64       `json:"-"`
	TagScore        float64       `json:"-"`
	TotalScore      float64       `json:"-"`
}

type PaginatedFeedResponseDTO = CursorBasedPaginationResponseDTO[*FeedPostResponseDTO]
