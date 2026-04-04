package handlers

import (
	"errors"
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
)

// Maps domain errors to HTTP responses
func handleError(w http.ResponseWriter, r *http.Request, err error, logger logging.Logger) {
	switch {
	case
		errors.Is(err, domain.ErrPostNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrCommentNotFound):
		respondWithError(w, http.StatusNotFound, err.Error(), logger)

	case errors.Is(err, domain.ErrPostVersionConflict):
		respondWithError(w, http.StatusConflict, err.Error(), logger)

	case errors.Is(err, domain.ErrUserAlreadyFollowing):
		respondWithError(w, http.StatusNoContent, err.Error(), logger)

	default:
		handleInternalServerError(w, r, err, logger)
	}
}

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error, logger logging.Logger) {
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
