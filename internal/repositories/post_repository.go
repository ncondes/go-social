package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()
	query := `
	INSERT INTO posts (title, content, user_id, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, version, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return handleDBError(err, resourcePost)
	}

	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	SELECT
		p.id,
		p.title,
		p.content,
		p.tags,
		p.user_id,
		p.version,
		p.created_at,
		p.updated_at,
		u.id,
		u.first_name,
		u.last_name,
		u.username,
		u.email,
		u.created_at,
		u.updated_at,
		COUNT(c.id) AS comment_count
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN comments c ON p.id = c.post_id
	WHERE p.id = $1
	GROUP BY p.id, u.id`

	row := r.db.QueryRowContext(ctx, query, postID)

	var result domain.PostWithDetails

	err := row.Scan(
		&result.Post.ID,
		&result.Post.Title,
		&result.Post.Content,
		pq.Array(&result.Post.Tags),
		&result.Post.UserID,
		&result.Post.Version,
		&result.Post.CreatedAt,
		&result.Post.UpdatedAt,
		&result.Author.ID,
		&result.Author.FirstName,
		&result.Author.LastName,
		&result.Author.Username,
		&result.Author.Email,
		&result.Author.CreatedAt,
		&result.Author.UpdatedAt,
		&result.CommentCount,
	)
	if err != nil {
		return nil, handleDBError(err, resourcePost)
	}

	return &result, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	setFields := []string{}
	args := []any{}
	argIndex := 1

	if post.Title != "" {
		setFields = append(setFields, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, post.Title)
		argIndex++
	}

	if post.Content != "" {
		setFields = append(setFields, fmt.Sprintf("content = $%d", argIndex))
		args = append(args, post.Content)
		argIndex++
	}

	if len(post.Tags) > 0 {
		setFields = append(setFields, fmt.Sprintf("tags = $%d", argIndex))
		args = append(args, pq.Array(post.Tags))
		argIndex++
	}

	if len(setFields) == 0 {
		return nil // No fields to update
	}

	setFields = append(setFields, "updated_at = NOW()", "version = version + 1")
	args = append(args, post.ID, post.Version)

	// Optimistic locking via version column
	query := fmt.Sprintf(
		"UPDATE posts SET %s WHERE id = $%d AND version = $%d RETURNING id, title, content, user_id, tags, version, created_at, updated_at",
		strings.Join(setFields, ", "),
		argIndex,
		argIndex+1,
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		args...,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.Version,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return handleDBError(err, resourcePost)
		}
		// Zero rows: distinguish between wrong ID and stale version
		var exists bool
		_ = r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", post.ID).Scan(&exists)
		if !exists {
			return domain.ErrPostNotFound
		}
		return domain.ErrPostVersionConflict
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, postID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
	DELETE FROM posts
	WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return handleDBError(err, resourcePost)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return handleDBError(err, resourcePost)
	}

	if rows == 0 {
		return domain.ErrPostNotFound
	}

	return nil
}
