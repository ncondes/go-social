package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/ncondes/go/social/internal/domain"
)

const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
	notNullViolation    = "23502"
)

type resourceType string

const (
	resourcePost    resourceType = "post"
	resourceUser    resourceType = "user"
	resourceComment resourceType = "comment"
)

func handleDBError(err error, resource resourceType) error {
	// Handle SQL errors
	if err == sql.ErrNoRows {
		switch resource {
		case resourcePost:
			return domain.ErrPostNotFound
		case resourceUser:
			return domain.ErrUserNotFound
		case resourceComment:
			return domain.ErrCommentNotFound
		default:
			return err
		}
	}
	// Handle PostgreSQL specific errors
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}

	// TODO: remove this
	fmt.Println(pqErr.Code, pqErr.Constraint)

	switch pqErr.Code {
	case foreignKeyViolation:
		return translateForeignKeyViolationError(pqErr)
	case uniqueViolation:
		return translateUniqueViolationError(pqErr)
	default:
		return pqErr
	}
}

func translateForeignKeyViolationError(pqErr *pq.Error) error {
	switch pqErr.Constraint {
	case "fk_posts_user_id":
		return domain.ErrUserNotFound
	case "fk_comments_post_id":
		return domain.ErrPostNotFound
	case "fk_comments_user_id":
		return domain.ErrUserNotFound
	case "fk_followers_user_id":
		return domain.ErrUserNotFound
	case "fk_followers_follower_id":
		return domain.ErrUserNotFound
	default:
		return pqErr
	}
}

func translateUniqueViolationError(pqErr *pq.Error) error {
	switch pqErr.Constraint {
	case "pk_followers":
		return domain.ErrUserAlreadyFollowing
	default:
		return pqErr
	}
}
