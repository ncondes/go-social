package repositories

import (
	"context"
	"database/sql"

	"github.com/ncondes/go/social/internal/domain"
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

func (r *CommentRepository) GetManyByPostID(ctx context.Context, postID int64) ([]*domain.CommentWithAuthor, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		c.id,
		c.post_id,
		c.user_id,
		c.content,
		c.created_at,
		c.updated_at,
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		u.email,
		u.created_at,
		u.updated_at
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = $1
	ORDER BY c.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, handleDBError(err, resourceComment)
	}

	defer rows.Close()

	var comments []*domain.CommentWithAuthor

	for rows.Next() {
		var result domain.CommentWithAuthor

		if err := rows.Scan(
			&result.Comment.ID,
			&result.Comment.PostID,
			&result.Comment.UserID,
			&result.Comment.Content,
			&result.Comment.CreatedAt,
			&result.Comment.UpdatedAt,
			&result.Author.ID,
			&result.Author.FirstName,
			&result.Author.LastName,
			&result.Author.Username,
			&result.Author.Email,
			&result.Author.CreatedAt,
			&result.Author.UpdatedAt,
		); err != nil {
			return nil, handleDBError(err, resourceComment)
		}

		comments = append(comments, &result)
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
