package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/services"
)

type FeedHandler struct {
	feedService *services.FeedService
}

func NewFeedHandler(feedService *services.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

func (h *FeedHandler) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	userID := int64(51) // TODO: get userID from auth middleware in the future

	options, err := ParseFeedQueryOptions(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.feedService.GetUserFeed(r.Context(), userID, options)
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	if err := respondWithPagination(w, http.StatusOK, result.Posts, dtos.CursorBasedPaginationMeta{
		Limit:      options.Pagination.Limit,
		NextCursor: result.NextCursor,
	}); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}
