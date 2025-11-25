package repositories

import (
	"database/sql"

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

	switch pqErr.Code {
	case foreignKeyViolation:
		return translateForeignKeyError(pqErr)
	default:
		return pqErr
	}
}

func translateForeignKeyError(pqErr *pq.Error) error {
	switch pqErr.Constraint {
	case "fk_post":
		return domain.ErrPostNotFound
	case "fk_user":
		return domain.ErrUserNotFound
	default:
		return pqErr
	}
}
