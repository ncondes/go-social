package domain

import (
	"context"
	"time"
)

type FeedPost struct {
	Post
	Author          User
	CommentCount    int
	RecencyScore    float64
	EngagementScore float64
	TagScore        float64
	TotalScore      float64
}

type FeedCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        int64     `json:"id"`
}

type FeedPaginationOptions struct {
	Limit  int
	Cursor *FeedCursor
}

type FeedFilterOptions struct {
	Since  *time.Time
	Until  *time.Time
	Search string
	Tags   []string
}

type FeedQueryOptions struct {
	Pagination FeedPaginationOptions
	Filters    FeedFilterOptions
}

type FeedRepositoryInterface interface {
	GetUserFeed(ctx context.Context, userID int64, options *FeedQueryOptions) ([]*FeedPost, error)
	GetUserTagInterests(ctx context.Context, userID int64) (map[string]int, error)
}

type FeedServiceInterface interface {
	GetUserFeed(ctx context.Context, userID int64, options *FeedQueryOptions) ([]*FeedPost, error)
}
