package domain

import (
	"context"
	"time"

	"github.com/ncondes/go/social/internal/dtos"
)

type FeedRepository interface {
	GetUserFeed(ctx context.Context, userID int64, options *FeedQueryOptions) ([]*dtos.FeedPostResponseDTO, error)
	GetUserTagInterests(ctx context.Context, userID int64) (map[string]int, error)
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
