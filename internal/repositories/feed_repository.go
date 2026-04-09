package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/ncondes/go/social/internal/domain"
)

type FeedRepository struct {
	db *sql.DB
}

func NewFeedRepository(db *sql.DB) *FeedRepository {
	return &FeedRepository{db: db}
}

func (r *FeedRepository) GetUserFeed(ctx context.Context, userID int64, options *domain.FeedQueryOptions) ([]*domain.FeedPost, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query, args := r.buildFeedQuery(userID, options)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, handleDBError(err, resourcePost)
	}

	defer rows.Close()

	var feeds []*domain.FeedPost

	for rows.Next() {
		var result domain.FeedPost
		var authorFirstName, authorLastName string

		if err := rows.Scan(
			&result.Post.ID,
			&result.Post.Title,
			&result.Post.Content,
			pq.Array(&result.Post.Tags),
			&result.Post.CreatedAt,
			&result.Post.UpdatedAt,
			&result.Author.ID,
			&result.Author.Username,
			&authorFirstName,
			&authorLastName,
			&result.CommentCount,
			&result.RecencyScore,
			&result.EngagementScore,
		); err != nil {
			return nil, handleDBError(err, resourcePost)
		}

		// Set author name fields
		result.Author.FirstName = authorFirstName
		result.Author.LastName = authorLastName

		feeds = append(feeds, &result)
	}

	return feeds, nil
}

func (r *FeedRepository) buildFeedQuery(userID int64, options *domain.FeedQueryOptions) (string, []any) {
	baseQuery := `
		SELECT
			p.id,
			p.title,
			p.content,
			p.tags,
			p.created_at,
			p.updated_at,
			u.id AS author_id,
			u.username,
			u.first_name,
			u.last_name,
			COUNT(c.id) AS comment_count,
			-- Recency score (0-1, decays over time)
			1.0 / (1.0 + EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 86400.0) as recency_score,
			-- Engagement score (0-1, normalized)
			LEAST(COUNT(DISTINCT c.id) / 10.0, 1.0) as engagement_score
		FROM posts p 
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id`

	var conditions []string
	args := []any{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("p.user_id IN (SELECT f.user_id FROM followers f WHERE f.follower_id = $%d)", argIndex))
	args = append(args, userID)
	argIndex++

	if options.Filters.Since != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
		args = append(args, *options.Filters.Since)
		argIndex++
	}

	if options.Filters.Until != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
		args = append(args, *options.Filters.Until)
		argIndex++
	}

	if options.Filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(p.title ILIKE $%d OR p.content ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+options.Filters.Search+"%")
		argIndex++
	}

	if len(options.Filters.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("p.tags && $%d", argIndex))
		args = append(args, pq.Array(options.Filters.Tags))
		argIndex++
	}

	if options.Pagination.Cursor != nil {
		conditions = append(conditions, fmt.Sprintf("(p.created_at, p.id) < ($%d, $%d)", argIndex, argIndex+1))
		args = append(args, options.Pagination.Cursor.CreatedAt, options.Pagination.Cursor.ID)
		argIndex += 2
	}

	if len(conditions) > 0 {
		baseQuery += "\n\t\tWHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(`
		GROUP BY p.id, u.id
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT $%d`, argIndex)
	args = append(args, options.Pagination.Limit)

	return baseQuery, args
}

func (r *FeedRepository) GetUserTagInterests(ctx context.Context, userID int64) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `
		SELECT DISTINCT unnest(p.tags::text[]) as tag, COUNT(*) as engagement_count
		FROM comments c
		JOIN posts p ON c.post_id = p.id
		WHERE c.user_id = $1
		GROUP BY tag
		ORDER BY engagement_count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, handleDBError(err, resourceFeed)
	}

	defer rows.Close()

	interests := make(map[string]int)

	for rows.Next() {
		var tag string
		var count int

		if err := rows.Scan(&tag, &count); err != nil {
			return nil, handleDBError(err, resourceFeed)
		}

		interests[tag] = count
	}

	return interests, nil
}
