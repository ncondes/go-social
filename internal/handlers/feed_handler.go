package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/packages/pagination"
)

type FeedHandler struct {
	feedService domain.FeedServiceInterface
}

func NewFeedHandler(feedService domain.FeedServiceInterface) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

func (h *FeedHandler) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	userID := int64(51) // TODO: get userID from auth middleware in the future

	options, err := h.parseFeedQueryOptions(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	feedPosts, err := h.feedService.GetUserFeed(r.Context(), userID, options)
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	responsePosts := make([]*dtos.FeedPostResponseDTO, len(feedPosts))
	for i, post := range feedPosts {
		responsePosts[i] = new(dtos.FeedPostResponseDTO).FromDomain(post)
	}

	nextCursor, err := h.buildNextCursor(feedPosts, options.Pagination.Limit)
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	response := dtos.FeedResponseDTO{
		Posts: responsePosts,
		Pagination: dtos.CursorBasedPaginationMeta{
			Limit:      options.Pagination.Limit,
			NextCursor: nextCursor,
		},
	}

	if err := respondWithData(w, http.StatusOK, response); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *FeedHandler) parseFeedQueryOptions(r *http.Request) (*domain.FeedQueryOptions, error) {
	pagination, err := h.parseFeedPaginationOptions(r)
	if err != nil {
		return nil, err
	}

	filters, err := h.parseFeedFilterOptions(r)
	if err != nil {
		return nil, err
	}

	return &domain.FeedQueryOptions{
		Pagination: pagination,
		Filters:    filters,
	}, nil
}

func (h *FeedHandler) parseFeedPaginationOptions(r *http.Request) (domain.FeedPaginationOptions, error) {
	limit, cursor, err := parseCursorPaginationParams[domain.FeedCursor](r)
	if err != nil {
		return domain.FeedPaginationOptions{}, err
	}

	return domain.FeedPaginationOptions{
		Limit:  limit,
		Cursor: cursor,
	}, nil
}

func (h *FeedHandler) parseFeedFilterOptions(r *http.Request) (domain.FeedFilterOptions, error) {
	filters := domain.FeedFilterOptions{}

	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return filters, errors.New("invalid since date format, use RFC3339 (e.g., 2024-11-26T00:00:00Z)")
		}
		filters.Since = &since
	}

	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		until, err := time.Parse(time.RFC3339, untilStr)
		if err != nil {
			return filters, errors.New("invalid until date format, use RFC3339 (e.g., 2024-11-26T23:59:59Z)")
		}
		filters.Until = &until
	}

	if filters.Since != nil && filters.Until != nil && filters.Since.After(*filters.Until) {
		return filters, errors.New("since date must be before until date")
	}

	filters.Search = strings.TrimSpace(r.URL.Query().Get("search"))
	filters.Tags = parseTagsParam(r, "tags")

	return filters, nil
}

func (h *FeedHandler) buildNextCursor(feedPosts []*domain.FeedPost, limit int) (string, error) {
	if len(feedPosts) < limit {
		return "", nil
	}

	lastPost := feedPosts[len(feedPosts)-1]
	return pagination.EncodeCursor(domain.FeedCursor{
		CreatedAt: lastPost.Post.CreatedAt,
		ID:        lastPost.Post.ID,
	})
}
