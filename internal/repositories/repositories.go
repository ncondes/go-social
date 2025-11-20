package repositories

import (
	"database/sql"

	"github.com/ncondes/go/social/internal/domain"
)

type Repositories struct {
	UserRepository domain.UserRepository
	PostRepository domain.PostRepository
}

func New(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepository: NewUserRepository(db),
		PostRepository: NewPostRepository(db),
	}
}
