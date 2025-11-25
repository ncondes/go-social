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
	UserRepository    domain.UserRepository
	PostRepository    domain.PostRepository
	CommentRepository domain.CommentRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository:    NewUserRepository(db),
		PostRepository:    NewPostRepository(db),
		CommentRepository: NewCommentRepository(db),
	}
}
