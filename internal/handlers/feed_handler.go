package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/packages/pagination"
)

type FeedHandler struct {
	feedService domain.FeedServiceInterface
	logger      logging.Logger
}

func NewFeedHandler(feedService domain.FeedServiceInterface, logger logging.Logger) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
		logger:      logger,
	}
}

func (h *FeedHandler) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	user := getAuthenticatedUserFromContext(r.Context())

	options, err := h.parseFeedQueryOptions(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	feedPosts, err := h.feedService.GetUserFeed(r.Context(), user.ID, options)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	responsePosts := make([]*dtos.FeedPostResponseDTO, len(feedPosts))
	for i, post := range feedPosts {
		responsePosts[i] = new(dtos.FeedPostResponseDTO).FromDomain(post)
	}

	nextCursor, err := h.buildNextCursor(feedPosts, options.Pagination.Limit)
	if err != nil {
		handleInternalServerError(w, r, err, h.logger)
		return
	}

	pagination := dtos.CursorBasedPaginationMetaDTO{
		Limit:      options.Pagination.Limit,
		NextCursor: nextCursor,
	}

	respondWithPaginatedData(w, http.StatusOK, responsePosts, pagination, h.logger)
}

func (h *FeedHandler) parseFeedQueryOptions(r *http.Request) (*domain.FeedQueryOptions, error) {
	limit, cursor, err := parseCursorPaginationParams[domain.FeedCursor](r)
	if err != nil {
		return nil, err
	}

	filters, err := h.parseFeedFilterOptions(r)
	if err != nil {
		return nil, err
	}

	return &domain.FeedQueryOptions{
		Pagination: domain.FeedPaginationOptions{
			Limit:  limit,
			Cursor: cursor,
		},
		Filters: filters,
	}, nil
}

func (h *FeedHandler) parseFeedFilterOptions(r *http.Request) (domain.FeedFilterOptions, error) {
	since, err := parseDateParam(r, "since")
	if err != nil {
		return domain.FeedFilterOptions{}, err
	}

	until, err := parseDateParam(r, "until")
	if err != nil {
		return domain.FeedFilterOptions{}, err
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))
	tags := parseTagsParam(r, "tags")

	isDateRangeInvalid := since != nil && until != nil && since.After(*until)
	if isDateRangeInvalid {
		return domain.FeedFilterOptions{}, errors.New("since date must be before until date")
	}

	return domain.FeedFilterOptions{
		Since:  since,
		Until:  until,
		Search: search,
		Tags:   tags,
	}, nil
}

func (h *FeedHandler) buildNextCursor(feedPosts []*domain.FeedPost, limit int) (string, error) {
	if len(feedPosts) < limit {
		return "", nil
	}

	lastPost := feedPosts[len(feedPosts)-1]
	return pagination.EncodeCursor(domain.FeedCursor{
		CreatedAt: lastPost.CreatedAt,
		ID:        lastPost.ID,
	})
}
