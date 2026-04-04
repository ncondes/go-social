package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/ncondes/go/social/internal/config"
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

func New(db *sql.DB, config *config.Config) *Repositories {
	return &Repositories{
		UserRepository:     NewUserRepository(db, config),
		PostRepository:     NewPostRepository(db),
		CommentRepository:  NewCommentRepository(db),
		FollowerRepository: NewFollowerRepository(db),
		FeedRepository:     NewFeedRepository(db),
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}
