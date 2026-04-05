package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/domain"
)

type UserRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewUserRepository(db *sql.DB, config *config.Config) *UserRepository {
	return &UserRepository{db: db, config: config}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO users (first_name, last_name, username, email, password)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return handleDBError(err, resourceUser)
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `SELECT id, first_name, last_name, username, email, password, is_active, created_at, updated_at
	FROM users
	WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, handleDBError(err, resourceUser)
	}

	return user, nil
}

func (r *UserRepository) CreateUserAndInvitation(
	ctx context.Context,
	user *domain.User,
	method string,
	token string,
) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		if err := r.createUserWithTx(ctx, tx, user); err != nil {
			return err
		}

		if err := r.createUserInvitationWithTx(ctx, tx, user.ID, method, token, r.config.MailConfig.Exp); err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) ActivateUser(ctx context.Context, token string) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
		defer cancel()
		// Find user that the token belongs to and check that the token is valid
		user, err := r.getUserByTokenWithTx(ctx, tx, token)
		if err != nil {
			return err
		}
		// Activate user
		active := true
		if err := r.updateUserWithTx(ctx, tx, user.ID, &domain.UserUpdate{IsActive: &active}); err != nil {
			return err
		}
		// Clean up the user's invitation
		if err := r.deleteUserInvitationWithTx(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int64) error {
	return withTx(r.db, ctx, func(tx *sql.Tx) error {
		if err := r.deleteUserWithTx(ctx, tx, userID); err != nil {
			return err
		}

		if err := r.deleteUserInvitationWithTx(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}

func (r *UserRepository) deleteUserWithTx(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return handleDBError(err, resourceUser)
	}

	return nil
}

func (r *UserRepository) deleteUserInvitationWithTx(ctx context.Context, tx *sql.Tx, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `DELETE FROM user_invitations WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return handleDBError(err, resourceUserInvitation)
	}

	return nil
}

func (r *UserRepository) createUserWithTx(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO users (first_name, last_name, username, email, password)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return handleDBError(err, resourceUser)
	}

	return nil
}

// TODO: we need to consider
// TODO: create a cron-job to clean up expired invitations
func (r *UserRepository) createUserInvitationWithTx(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	method string,
	token string,
	exp time.Duration, // expiration time
) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `INSERT INTO user_invitations (user_id, method, token_hash, created_at, expires_at)
	VALUES ($1, $2, $3, $4, $5)`

	expiresAt := time.Now().Add(exp)

	_, err := tx.ExecContext(ctx, query, userID, method, token, time.Now(), expiresAt)
	if err != nil {
		return handleDBError(err, resourceUserInvitation)
	}

	return nil
}

func (r *UserRepository) getUserByTokenWithTx(
	ctx context.Context,
	tx *sql.Tx,
	token string,
) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	query := `SELECT
		u.id,
		u.username,
		u.email,
		u.created_at,
		u.is_active
	FROM users u
	JOIN user_invitations ui ON u.id = ui.user_id
	WHERE ui.token_hash = $1 AND ui.expires_at > NOW()`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &domain.User{}
	err := tx.QueryRowContext(ctx, query, hashToken).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		return nil, handleDBError(err, resourceUser)
	}

	return user, nil
}

func (r *UserRepository) updateUserWithTx(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	update *domain.UserUpdate,
) error {
	if update == nil {
		return domain.ErrNoUserUpdateFields
	}

	query, args, err := buildUserUpdateQuery(userID, update)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return handleDBError(err, resourceUser)
	}
	return nil
}

func buildUserUpdateQuery(userID int64, u *domain.UserUpdate) (string, []any, error) {
	var sets []string
	var args []any
	n := 1

	add := func(column string, value any) {
		sets = append(sets, fmt.Sprintf("%s = $%d", column, n))
		args = append(args, value)
		n++
	}

	if u.FirstName != nil {
		add("first_name", *u.FirstName)
	}
	if u.LastName != nil {
		add("last_name", *u.LastName)
	}
	if u.Username != nil {
		add("username", *u.Username)
	}
	if u.Email != nil {
		add("email", *u.Email)
	}
	if u.IsActive != nil {
		add("is_active", *u.IsActive)
	}

	if len(sets) == 0 {
		return "", nil, domain.ErrNoUserUpdateFields
	}

	sets = append(sets, "updated_at = NOW()")
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(sets, ", "),
		n,
	)
	args = append(args, userID)
	return query, args, nil
}
