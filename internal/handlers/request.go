package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ncondes/go/social/packages/pagination"
)

const (
	defaultLimit = 20
	maxLimit     = 100
	minLimit     = 1
)

func parseCursorPaginationParams[T any](r *http.Request) (limit int, cursor *T, err error) {
	limit = defaultLimit

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		l, parseErr := strconv.Atoi(limitStr)
		if parseErr != nil {
			return defaultLimit, nil, fmt.Errorf("limit parameter must be a number")
		}
		if l < minLimit || l > maxLimit {
			return defaultLimit, nil, fmt.Errorf("limit parameter must be between %d and %d", minLimit, maxLimit)
		}

		limit = l
	}

	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		decoded, err := pagination.DecodeCursor[T](cursorStr)
		if err != nil {
			return defaultLimit, nil, fmt.Errorf("invalid cursor parameter")
		}

		cursor = &decoded
	}

	return limit, cursor, nil
}

func parseTagsParam(r *http.Request, key string) []string {
	tags := []string{}

	for _, value := range r.URL.Query()[key] {
		// Each value might be comma-separated
		for _, tag := range strings.Split(value, ",") {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}

	return tags
}

func parseDateParam(r *http.Request, key string) (*time.Time, error) {
	str := r.URL.Query().Get(key)
	if str == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, fmt.Errorf("invalid %s date format, use RFC3339 (e.g., 2024-11-26T00:00:00Z)", key)
	}

	return &parsed, nil
}
