package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func PostIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postIDParam := chi.URLParam(r, "postID")
		postID, err := strconv.ParseInt(postIDParam, 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid post ID")
			return
		}

		ctx := context.WithValue(r.Context(), "postID", postID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetPostIDFromContext(ctx context.Context) int64 {
	postID, _ := ctx.Value("postID").(int64)
	return postID
}
