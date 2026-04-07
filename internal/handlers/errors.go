package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
)

func handleError(w http.ResponseWriter, r *http.Request, err error, logger logging.Logger) {
	if errors.Is(err, context.Canceled) {
		return
	}
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		respondWithError(w, http.StatusUnauthorized, err.Error(), logger)
	case
		errors.Is(err, domain.ErrPostNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrCommentNotFound):
		respondWithError(w, http.StatusNotFound, err.Error(), logger)

	case
		errors.Is(err, domain.ErrPostVersionConflict),
		errors.Is(err, domain.ErrUserEmailTaken),
		errors.Is(err, domain.ErrUserUsernameTaken):
		respondWithError(w, http.StatusConflict, err.Error(), logger)

	case errors.Is(err, domain.ErrUserAlreadyFollowing):
		respondWithError(w, http.StatusNoContent, err.Error(), logger)

	default:
		handleInternalServerError(w, r, err, logger)
	}
}

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error, logger logging.Logger) {
	if err != nil && errors.Is(err, context.Canceled) {
		return
	}
	switch {
	case r != nil && err != nil:
		logger.Errorw("internal server error",
			"error", err,
			"method", r.Method,
			"path", r.URL.Path,
		)
	case r != nil:
		logger.Errorw("internal server error",
			"method", r.Method,
			"path", r.URL.Path,
		)
	case err != nil:
		logger.Errorw("internal server error",
			"error", err,
			"source", "json_response",
		)
	default:
		logger.Errorw("internal server error", "source", "json_response")
	}

	respondWithError(w, http.StatusInternalServerError, "internal server error", logger)
}
