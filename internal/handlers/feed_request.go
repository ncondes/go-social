package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

func ParseFeedPaginationOptions(r *http.Request) (domain.FeedPaginationOptions, error) {
	limit, cursor, err := ParseCursorPaginationParams[domain.FeedCursor](r)
	if err != nil {
		return domain.FeedPaginationOptions{}, err
	}

	return domain.FeedPaginationOptions{
		Limit:  limit,
		Cursor: cursor,
	}, nil
}

func ParseFeedFilterOptions(r *http.Request) (domain.FeedFilterOptions, error) {
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

func ParseFeedQueryOptions(r *http.Request) (*domain.FeedQueryOptions, error) {
	pagination, err := ParseFeedPaginationOptions(r)
	if err != nil {
		return nil, err
	}

	filters, err := ParseFeedFilterOptions(r)
	if err != nil {
		return nil, err
	}

	return &domain.FeedQueryOptions{
		Pagination: pagination,
		Filters:    filters,
	}, nil
}

func parseTagsParam(r *http.Request, key string) []string {
	var tags []string

	if tagsStr := r.URL.Query().Get(key); tagsStr != "" {
		// Split by comma and trim spaces ?tags=golang,docker
		for _, tag := range strings.Split(tagsStr, ",") {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	} else if tagsList := r.URL.Query()[key]; len(tagsList) > 0 {
		// Support multiple ?tags=golang&tags=docker
		for _, tag := range tagsList {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	return tags
}
