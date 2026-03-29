package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ncondes/go/social/packages/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCursorPaginationParams(t *testing.T) {
	t.Run("should return default limit when no params provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		limit, cursor, err := parseCursorPaginationParams[int](req)

		assert.NoError(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should parse valid limit when provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?limit=10", nil)

		limit, cursor, err := parseCursorPaginationParams[int](req)

		assert.NoError(t, err)
		assert.Equal(t, 10, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should return an error when limit is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?limit=invalid", nil)

		limit, cursor, err := parseCursorPaginationParams[int](req)

		assert.Error(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should return an error for limit below min", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?limit=-1", nil)

		limit, cursor, err := parseCursorPaginationParams[int](req)

		assert.Error(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should return an error for limit above max", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?limit=1001", nil)

		limit, cursor, err := parseCursorPaginationParams[int](req)

		assert.Error(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should parse a valid cursor", func(t *testing.T) {
		cursorValue := int64(12345)
		encoded, err := pagination.EncodeCursor(cursorValue)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test?cursor="+encoded, nil)

		limit, cursor, err := parseCursorPaginationParams[int64](req)

		assert.NoError(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Equal(t, cursorValue, *cursor)
	})

	t.Run("should return an error for invalid cursor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?cursor=invalid", nil)

		limit, cursor, err := parseCursorPaginationParams[int64](req)

		assert.Error(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})

	t.Run("should parse both limit and cursor", func(t *testing.T) {
		cursorValue := int64(12345)
		encodedCursor, err := pagination.EncodeCursor(cursorValue)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/test?limit=15&cursor="+encodedCursor, nil)

		limit, cursor, err := parseCursorPaginationParams[int64](req)

		assert.NoError(t, err)
		assert.Equal(t, 15, limit)
		assert.Equal(t, cursorValue, *cursor)
	})

	t.Run("should handle empty cursor", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?cursor=", nil)

		limit, cursor, err := parseCursorPaginationParams[int64](req)

		assert.NoError(t, err)
		assert.Equal(t, defaultLimit, limit)
		assert.Nil(t, cursor)
	})
}

func TestParseTagsParam(t *testing.T) {
	t.Run("should return empty slice when no tags provided", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		tags := parseTagsParam(req, "tags")

		assert.Empty(t, tags)
	})

	t.Run("should parse tags separated by comma", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?tags=golang,docker", nil)

		tags := parseTagsParam(req, "tags")

		assert.Equal(t, []string{"golang", "docker"}, tags)
	})

	t.Run("should parse tags from multiple query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?", nil)
		query := req.URL.Query()
		query.Add("tags", "golang")
		query.Add("tags", "docker")
		req.URL.RawQuery = query.Encode()

		tags := parseTagsParam(req, "tags")

		assert.Equal(t, []string{"golang", "docker"}, tags)
	})

	t.Run("should trim whitespace from tags", func(t *testing.T) {
		encodedSpace := "%20"
		req := httptest.NewRequest(http.MethodGet, "/test?tags=golang,"+encodedSpace+"docker,"+encodedSpace+"rust", nil)

		tags := parseTagsParam(req, "tags")

		assert.Equal(t, []string{"golang", "docker", "rust"}, tags)
	})

	t.Run("filters out empty tags from comma-separated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test?tags=golang,,docker", nil)

		tags := parseTagsParam(req, "tags")

		assert.Equal(t, []string{"golang", "docker"}, tags)
	})

	t.Run("filters out empty tags from multiple params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		q := req.URL.Query()
		q.Add("tags", "golang")
		q.Add("tags", "")
		q.Add("tags", "docker")
		req.URL.RawQuery = q.Encode()

		tags := parseTagsParam(req, "tags")

		assert.Equal(t, []string{"golang", "docker"}, tags)
	})
}
