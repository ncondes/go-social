# Cursor-Based Pagination

## Overview

The feed API uses cursor-based pagination for efficient, stable scrolling through posts. This prevents duplicates and skipped items when new content is added while users are browsing.

## How It Works

### First Request (No Cursor)
```bash
GET /v1/feed?limit=20
```

**Response:**
```json
{
  "data": [
    {
      "id": 187,
      "title": "Building Microservices",
      "content": "Here's what I learned...",
      "tags": ["golang", "architecture"],
      "author": {
        "id": 42,
        "username": "johndoe",
        "fullname": "John Doe"
      },
      "created_at": "2024-11-26T14:30:00Z",
      "updated_at": "2024-11-26T14:30:00Z",
      "comment_count": 15
    },
    {
      "id": 186,
      "title": "Docker Best Practices",
      "content": "Tips for containerization...",
      "tags": ["docker", "devops"],
      "author": {
        "id": 23,
        "username": "janedoe",
        "fullname": "Jane Doe"
      },
      "created_at": "2024-11-26T13:45:00Z",
      "updated_at": "2024-11-26T13:45:00Z",
      "comment_count": 8
    }
  ],
  "pagination": {
    "limit": 20,
    "next_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNlQxMzowMDowMFoiLCJpZCI6MTY1fQ=="
  }
}
```

### Next Page (With Cursor)
```bash
GET /v1/feed?limit=20&cursor=eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNlQxMzowMDowMFoiLCJpZCI6MTY1fQ==
```

**Response:**
```json
{
  "data": [
    {
      "id": 164,
      "title": "GraphQL vs REST",
      "content": "Comparing API approaches...",
      "tags": ["api", "webdev"],
      "author": {
        "id": 15,
        "username": "alexsmith",
        "fullname": "Alex Smith"
      },
      "created_at": "2024-11-26T12:30:00Z",
      "updated_at": "2024-11-26T12:30:00Z",
      "comment_count": 12
    }
  ],
  "pagination": {
    "limit": 20,
    "next_cursor": "eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNlQxMTowMDowMFoiLCJpZCI6MTQ1fQ=="
  }
}
```

### Last Page (No More Data)
```bash
GET /v1/feed?limit=20&cursor=eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNVQwOTowMDowMFoiLCJpZCI6MjB9
```

**Response:**
```json
{
  "data": [
    {
      "id": 5,
      "title": "Getting Started with Go",
      "content": "A beginner's guide...",
      "tags": ["golang", "tutorial"],
      "author": {
        "id": 1,
        "username": "admin",
        "fullname": "Admin User"
      },
      "created_at": "2024-11-20T10:00:00Z",
      "updated_at": "2024-11-20T10:00:00Z",
      "comment_count": 25
    }
  ],
  "pagination": {
    "limit": 20,
    "next_cursor": ""
  }
}
```

## How Cursors Are Constructed

### Cursor Structure

Cursors are **base64-encoded JSON** objects containing the position of the last item in the current page. This allows the server to efficiently fetch the next set of results.

**Cursor Components:**
```json
{
  "created_at": "2024-11-26T13:45:00Z",
  "id": 186
}
```

**Encoded Cursor:**
```
eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNlQxMzo0NTowMFoiLCJpZCI6MTg2fQ==
```

### Why This Structure?

1. **Composite Key** - Uses `(created_at, id)` tuple for stable ordering
   - `created_at` - Primary sort field (newest first)
   - `id` - Tiebreaker for posts with identical timestamps

2. **SQL Efficiency** - Translates to an indexed WHERE clause:
   ```sql
   WHERE (created_at, id) < ($cursor.created_at, $cursor.id)
   ORDER BY created_at DESC, id DESC
   ```

3. **Stability** - Even if new posts are added, the cursor always points to the same position

### Cursor Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Server returns posts + builds cursor from last post     │
│    Last post: {id: 186, created_at: "2024-11-26T13:45:00Z"} │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ 2. Server encodes cursor to base64                          │
│    {"created_at": "...", "id": 186}                         │
│    → eyJjcmVhdGVkX2F0IjoiMjAyNC0xMS0yNlQxMzo0NTowMFoiLCJpZCI6MTg2fQ== │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ 3. Client receives cursor in response                       │
│    {"pagination": {"next_cursor": "eyJjcmVh..."}}           │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ 4. Client sends cursor in next request                      │
│    GET /v1/feed?cursor=eyJjcmVh...                          │
└─────────────────────────────────────────────────────────────┘
                           ↓
┌─────────────────────────────────────────────────────────────┐
│ 5. Server decodes cursor and fetches posts AFTER position   │
│    WHERE (created_at, id) < ("2024-11-26T13:45:00Z", 186)   │
└─────────────────────────────────────────────────────────────┘
```

### Example: Decoding a Cursor

**Encoded:**
```
eyJjcmVhdGVkX2F0IjoiMjAyNS0xMS0yNVQyMzoxMToyNFoiLCJpZCI6MTk4fQ==
```

**Decoded (base64 → JSON):**
```json
{
  "created_at": "2025-11-25T23:11:24Z",
  "id": 198
}
```

**Meaning:** "Fetch posts created before 2025-11-25T23:11:24Z, or if same timestamp, with ID < 198"

## Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | integer | No | 20 | Number of posts per page (max: 100) |
| `cursor` | string | No | - | Base64-encoded pagination cursor from previous response |

## Client Implementation

### JavaScript Example
```javascript
class FeedClient {
  constructor() {
    this.posts = [];
    this.nextCursor = null;
  }

  async loadInitial() {
    const response = await fetch('/v1/feed?limit=20');
    const data = await response.json();
    
    this.posts = data.data;
    this.nextCursor = data.pagination.next_cursor;
    
    return this.posts;
  }

  async loadMore() {
    if (!this.nextCursor) {
      return []; // No more data
    }

    const response = await fetch(`/v1/feed?limit=20&cursor=${this.nextCursor}`);
    const data = await response.json();
    
    this.posts = [...this.posts, ...data.data];
    this.nextCursor = data.pagination.next_cursor;
    
    return data.data;
  }

  hasMore() {
    return this.nextCursor !== '';
  }

  async refresh() {
    return this.loadInitial();
  }
}

// Usage
const feed = new FeedClient();

// Initial load
await feed.loadInitial();

// Infinite scroll
window.addEventListener('scroll', async () => {
  if (isNearBottom() && feed.hasMore()) {
    await feed.loadMore();
  }
});

// Pull to refresh
await feed.refresh();
```

### Go Client Example
```go
type FeedClient struct {
    baseURL    string
    httpClient *http.Client
    nextCursor string
}

func (c *FeedClient) LoadFeed(limit int) (*FeedResponse, error) {
    url := fmt.Sprintf("%s/v1/feed?limit=%d", c.baseURL, limit)
    if c.nextCursor != "" {
        url += fmt.Sprintf("&cursor=%s", c.nextCursor)
    }

    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var feedResp FeedResponse
    if err := json.NewDecoder(resp.Body).Decode(&feedResp); err != nil {
        return nil, err
    }

    c.nextCursor = feedResp.Pagination.NextCursor
    return &feedResp, nil
}

func (c *FeedClient) HasMore() bool {
    return c.nextCursor != ""
}
```

## Benefits

✅ **No Duplicates** - New posts don't shift pagination  
✅ **Consistent** - Same results when scrolling, even with new content  
✅ **Performant** - Uses indexed WHERE clause instead of OFFSET  
✅ **Real-time Friendly** - Works well with live updates  

## Refresh vs Load More

| Action | Cursor | Result |
|--------|--------|--------|
| **Load More** | Use `next_cursor` | Older posts (scroll down) |
| **Refresh** | No cursor | Latest posts (reload feed) |

## Technical Details

### Server-Side Implementation

**Building the Cursor (Service Layer):**
```go
func (s *FeedService) buildNextCursor(feed []*FeedPost, limit int) (string, error) {
    if len(feed) < limit {
        return "", nil // No more data
    }

    lastPost := feed[len(feed)-1]
    return pagination.EncodeCursor(domain.FeedCursor{
        CreatedAt: lastPost.CreatedAt,
        ID:        lastPost.ID,
    })
}
```

**Using the Cursor (Repository Layer):**
```sql
-- First page (no cursor)
SELECT * FROM posts
WHERE user_id IN (SELECT user_id FROM followers WHERE follower_id = $1)
ORDER BY created_at DESC, id DESC
LIMIT $2

-- Subsequent pages (with cursor)
SELECT * FROM posts
WHERE user_id IN (SELECT user_id FROM followers WHERE follower_id = $1)
  AND (created_at, id) < ($2, $3)  -- Cursor condition
ORDER BY created_at DESC, id DESC
LIMIT $4
```

### Database Indexes

For optimal performance, ensure you have a composite index:
```sql
CREATE INDEX idx_posts_created_at_id ON posts(created_at DESC, id DESC);
```

## Important Notes

- ⚠️ **Don't parse cursors** - They are opaque tokens; structure may change
- ⚠️ **Don't modify cursors** - Always use them exactly as received
- ✅ **Empty cursor = end of data** - `next_cursor: ""` means no more results
- ✅ **Refresh pattern** - Omit cursor to get latest posts (pull-to-refresh)
- ✅ **Maximum limit** - Server enforces max 100 posts per request
- ✅ **Stable ordering** - Composite key ensures consistent pagination even with concurrent inserts
