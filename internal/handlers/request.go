package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ncondes/go/social/packages/pagination"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func parseCursorPaginationParams[T any](r *http.Request) (limit int, cursor *T, err error) {
	limit = DefaultLimit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 && l <= MaxLimit {
			limit = l
		}
	}

	if cursorStr := r.URL.Query().Get("cursor"); cursorStr != "" {
		decoded, err := pagination.DecodeCursor[T](cursorStr)
		if err != nil {
			return 0, nil, err
		}
		cursor = &decoded
	}

	return limit, cursor, nil
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
