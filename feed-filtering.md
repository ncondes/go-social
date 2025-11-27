# Feed Filtering

## Overview

The feed API supports multiple filtering options to help users find relevant content. All filters can be combined and work seamlessly with cursor-based pagination.

## Available Filters

### Tags
Filter posts by one or more tags (OR logic - posts with ANY of the specified tags).

```bash
# Single tag
GET /v1/feed?tags=golang

# Multiple tags (comma-separated)
GET /v1/feed?tags=golang,docker,kubernetes

# Multiple tags (multiple params)
GET /v1/feed?tags=golang&tags=docker
```

### Date Range
Filter posts by creation date using `since` (start date) and/or `until` (end date).

**Date Format:** RFC3339 (e.g., `2024-11-26T00:00:00Z`)

```bash
# Posts since a specific date
GET /v1/feed?since=2024-11-01T00:00:00Z

# Posts until a specific date
GET /v1/feed?until=2024-11-30T23:59:59Z

# Posts within a date range
GET /v1/feed?since=2024-11-01T00:00:00Z&until=2024-11-30T23:59:59Z
```

### Search
Search for posts containing specific text in title or content (case-insensitive).

```bash
# Search for "microservices"
GET /v1/feed?search=microservices

# Search with spaces (URL-encoded)
GET /v1/feed?search=docker%20best%20practices
```

## Combining Filters

All filters can be combined for powerful queries:

```bash
# Golang posts from last month with "docker" in content
GET /v1/feed?tags=golang&since=2024-10-01T00:00:00Z&until=2024-10-31T23:59:59Z&search=docker

# Recent posts about Kubernetes
GET /v1/feed?tags=kubernetes&since=2024-11-20T00:00:00Z

# Search with pagination
GET /v1/feed?search=microservices&limit=10&cursor=eyJjcmVh...
```

## Query Parameters Reference

| Parameter | Type | Required | Description | Example |
|-----------|------|----------|-------------|---------|
| `limit` | integer | No | Posts per page (default: 20, max: 100) | `20` |
| `cursor` | string | No | Pagination cursor from previous response | `eyJjcmVh...` |
| `tags` | string | No | Comma-separated tags or multiple params | `golang,docker` |
| `since` | string (RFC3339) | No | Start date (inclusive) | `2024-11-01T00:00:00Z` |
| `until` | string (RFC3339) | No | End date (inclusive) | `2024-11-30T23:59:59Z` |
| `search` | string | No | Search query for title/content | `microservices` |

## Architecture

### Request Parsing Structure

The feed filtering implementation follows clean architecture with clear separation of concerns:

```
internal/handlers/
├── request.go          # Generic utilities
│   └── ParseCursorPaginationParams[T]
│
├── feed_request.go     # Feed-specific parsing
│   ├── ParseFeedPaginationOptions()
│   ├── ParseFeedFilterOptions()
│   └── ParseFeedQueryOptions()
│
└── feed_handler.go     # HTTP handlers
    └── GetUserFeed()
```

### Domain Model

```go
// Pagination options
type FeedPaginationOptions struct {
    Limit  int
    Cursor *FeedCursor
}

// Filter options
type FeedFilterOptions struct {
    Since  *time.Time
    Until  *time.Time
    Search string
    Tags   []string
}

// Complete query options
type FeedQueryOptions struct {
    Pagination FeedPaginationOptions
    Filters    FeedFilterOptions
}
```

### Benefits

- **Single Responsibility** - Each function has one clear purpose
- **Composable** - Can use pagination without filters, or vice versa
- **Testable** - Easy to unit test each function independently
- **Reusable** - Generic utilities work for any resource
- **Maintainable** - Changes to feed parsing don't affect generic code

## Database Performance

### Required Indexes

```sql
-- Composite index for cursor pagination
CREATE INDEX idx_posts_created_at_id ON posts(created_at DESC, id DESC);

-- GIN index for tag filtering
CREATE INDEX idx_posts_tags ON posts USING GIN(tags);

-- Index for date range queries
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

### Query Efficiency

- **Tag filtering** uses PostgreSQL's `&&` (array overlap) operator - very efficient with GIN index
- **Date filtering** uses simple comparison operators - efficient with B-tree index
- **Search** uses `ILIKE` - consider full-text search for large datasets
- **Dynamic query building** only adds WHERE clauses for provided filters

## Example Response

```json
{
  "data": [
    {
      "id": 199,
      "title": "Understanding Microservices",
      "content": "Best practices from industry experts...",
      "tags": ["golang", "cloud", "devops"],
      "author": {
        "id": 99,
        "fullname": "David Moore",
        "username": "dmoore48"
      },
      "created_at": "2025-11-25T23:11:24Z",
      "updated_at": "2025-11-25T23:11:24Z",
      "comment_count": 2
    }
  ],
  "pagination": {
    "limit": 20,
    "next_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNS0xMS0yNVQyMzoxMToyNFoiLCJpZCI6MTk5fQ=="
  }
}
```

## Best Practices

✅ **Use specific filters** - Reduce result set for better performance  
✅ **Combine filters** - Create powerful, targeted queries  
✅ **Keep filters consistent** - Use same filters across paginated requests  
✅ **Use RFC3339 dates** - Standard format for date/time  
✅ **URL-encode search** - Properly encode special characters  

❌ **Don't change filters mid-pagination** - Cursor becomes invalid  
❌ **Don't use huge date ranges** - Can impact performance  
❌ **Don't skip validation** - Always validate user input
