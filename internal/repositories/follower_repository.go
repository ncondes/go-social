package repositories

import (
	"context"
	"database/sql"
)

type FollowerRepository struct {
	db *sql.DB
}

func NewFollowerRepository(db *sql.DB) *FollowerRepository {
	return &FollowerRepository{db: db}
}

func (r *FollowerRepository) FollowUser(ctx context.Context, userID int64, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO followers (user_id, follower_id)
	VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return handleDBError(err, resourceUser)
	}

	return nil
}

func (r *FollowerRepository) UnfollowUser(ctx context.Context, userID int64, followerID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2`

	_, err := r.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return handleDBError(err, resourceUser)
	}

	return nil
}
