package repositories

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/ncondes/go/social/internal/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `INSERT INTO posts (title, content, user_id, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}
