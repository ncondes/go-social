package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
)

// Maps domain errors to HTTP responses
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case
		errors.Is(err, domain.ErrPostNotFound),
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrCommentNotFound):
		respondWithError(w, http.StatusNotFound, err.Error())

	case errors.Is(err, domain.ErrUserAlreadyFollowing):
		respondWithError(w, http.StatusNoContent, err.Error())

	default:
		handleInternalServerError(w, r, err)
	}
}

func handleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: replace with proper logging
	log.Printf("internal servererror: %s, method: %s, path: %s", err.Error(), r.Method, r.URL.Path)

	respondWithError(w, http.StatusInternalServerError, "internal server error")
}
