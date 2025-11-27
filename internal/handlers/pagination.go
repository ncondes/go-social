package handlers

import (
	"net/http"
	"strconv"

	"github.com/ncondes/go/social/packages/pagination"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

func ParseCursorPaginationParams[T any](r *http.Request) (limit int, cursor *T, err error) {
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
