package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/ncondes/go/social/internal/logging"
)

func PostIDMiddleware(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			postIDParam := chi.URLParam(r, "postID")
			postID, err := strconv.ParseInt(postIDParam, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid post ID", logger)
				return
			}

			ctx := context.WithValue(r.Context(), postIDContextKey, postID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDMiddleware(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDParam := chi.URLParam(r, "userID")
			userID, err := strconv.ParseInt(userIDParam, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid user ID", logger)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
