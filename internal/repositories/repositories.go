package repositories

import (
	"database/sql"
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

const (
	queryTimeoutDuration = 5 * time.Second
)

type Repositories struct {
	UserRepository     domain.UserRepositoryInterface
	PostRepository     domain.PostRepositoryInterface
	CommentRepository  domain.CommentRepositoryInterface
	FollowerRepository domain.FollowerRepositoryInterface
	FeedRepository     domain.FeedRepositoryInterface
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(db),
		PostRepository:     NewPostRepository(db),
		CommentRepository:  NewCommentRepository(db),
		FollowerRepository: NewFollowerRepository(db),
		FeedRepository:     NewFeedRepository(db),
	}
}
