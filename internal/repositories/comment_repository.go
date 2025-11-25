package repositories

import (
	"context"
	"database/sql"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	INSERT INTO comments (post_id, user_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return handleDBError(err, resourceComment)
	}

	return nil
}

func (r *CommentRepository) GetManyByPostID(ctx context.Context, postID int64) ([]*dtos.CommentResponseDTO, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		c.id,
		c.post_id,
		c.content,
		c.created_at,
		c.updated_at,
		u.id,
		u.username
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, handleDBError(err, resourceComment)
	}

	defer rows.Close()

	var comments []*dtos.CommentResponseDTO

	for rows.Next() {
		var comment dtos.CommentResponseDTO

		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Author.ID,
			&comment.Author.Username,
		); err != nil {
			return nil, handleDBError(err, resourceComment)
		}

		comments = append(comments, &comment)
	}

	return comments, nil
}

func (r *CommentRepository) GetCountByPostID(ctx context.Context, postID int64) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	SELECT COUNT(*)
	FROM comments
	WHERE post_id = $1`

	row := r.db.QueryRowContext(ctx, query, postID)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return 0, handleDBError(err, resourceComment)
	}

	return count, nil
}
