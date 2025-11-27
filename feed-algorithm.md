# Feed Ranking Algorithm

## Overview

The feed uses a **ranked algorithm** that balances three key factors:
1. **Recency** - Newer posts are prioritized
2. **Engagement** - Posts with more interaction (comments) rank higher  
3. **Tag Interest** - Posts matching user's interests (inferred from comment history) get boosted

This creates a personalized, relevant feed that surfaces both fresh and popular content aligned with user preferences.

---

## The Formula

```
final_score = (recency_score * 0.4) + (engagement_score * 0.3) + (tag_interest_score * 0.3)
```

Where:

- **recency_score**: 0-1 (normalized, decays over time)
- **engagement_score**: 0-1 (normalized by comment count)
- **tag_interest_score**: 0-1 (based on user's tag preferences)

---

## Components

### Part 1: Recency Score

```
recency_score = 1.0 / (1.0 + age_in_days)
```

**What it does:** Gives higher scores to newer posts, with exponential decay over time.

**How it works:**

- Score ranges from 0 to 1
- Brand new posts get score close to 1.0
- Score decreases as post ages
- Uses days as the time unit

**Examples:**

- Post from **1 hour ago** (0.04 days): `1.0 / (1.0 + 0.04) = 0.96`
- Post from **12 hours ago** (0.5 days): `1.0 / (1.0 + 0.5) = 0.67`
- Post from **1 day ago**: `1.0 / (1.0 + 1.0) = 0.50`
- Post from **7 days ago**: `1.0 / (1.0 + 7.0) = 0.125`

> Newer posts have scores closer to 1.0, older posts decay toward 0.

---

### Part 2: Engagement Score

```
engagement_score = min(comment_count / 10.0, 1.0)
```

**What it does:** Gives higher scores to posts with more comments, normalized to 0-1 range.

**How it works:**

- Score ranges from 0 to 1
- 10 comments = maximum score of 1.0
- More than 10 comments still capped at 1.0
- Linear scaling up to the cap

**Examples:**

- **0 comments**: `0 / 10.0 = 0.0`
- **1 comment**: `1 / 10.0 = 0.1`
- **5 comments**: `5 / 10.0 = 0.5`
- **10 comments**: `10 / 10.0 = 1.0`
- **15 comments**: `min(15 / 10.0, 1.0) = 1.0`

> Posts with more engagement get higher scores, capped at 1.0 for fairness.

---

### Part 3: Tag Interest Score (Personalization)

**What it does:** Boosts posts containing tags the user has previously engaged with through comments.

**How it works:**

1. **Track User Interests** - Analyze which tags appear in posts the user has commented on
2. **Calculate Tag Scores** - Match post tags against user interests
3. **Apply Weights** - Exact matches get full weight, related tags get partial weight

**Algorithm:**

```go
// For each post tag:
for postTag in post.Tags {
    // Exact match: full weight
    if userInterests[postTag] exists {
        score += normalize(userInterests[postTag])
    }
    
    // Related tags: 50% weight
    for userTag in userInterests {
        if areRelated(postTag, userTag) {
            score += normalize(userInterests[userTag]) * 0.5
        }
    }
}

// Normalize by tag count to avoid bias
tagScore = score / len(post.Tags)
```

**Related Tags Map:**

```go
relatedTags = {
    "store":     ["shop", "retail", "ecommerce", "shopping"],
    "tech":      ["technology", "programming", "software", "coding"],
    "food":      ["cooking", "recipe", "restaurant", "cuisine"],
    "travel":    ["vacation", "tourism", "adventure", "explore"],
    "fitness":   ["health", "workout", "exercise", "gym"],
}
```

**Examples:**

**User has commented on:**

- Posts with "tech" tag: 5 times
- Posts with "programming" tag: 3 times
- Posts with "food" tag: 2 times

**Post A:** Tags: ["tech", "startup"]

```
- "tech" exact match: 5/5 = 1.0
- "startup" no match: 0
- Score: 1.0 / 2 = 0.5
```

**Post B:** Tags: ["software", "tutorial"]

```
- "software" related to "tech": (5/5) * 0.5 = 0.5
- "software" related to "programming": (3/5) * 0.5 = 0.3
- "tutorial" no match: 0
- Score: (0.5 + 0.3) / 2 = 0.4
```

**Post C:** Tags: ["travel", "photography"]

```
- No matches
- Score: 0.0
```

> Posts matching user interests get boosted, creating a personalized feed.

---

## Combined Score (Enhanced Version)

```
final_score = (recency_score * 0.4) + (engagement_score * 0.3) + (tag_interest_score * 0.3)
```

### Score Normalization

All three components are normalized to 0-1 range:

**Recency Score:**

```
recency_score = 1.0 / (1.0 + age_in_days)
```

**Engagement Score:**

```
engagement_score = min(comment_count / 10.0, 1.0)
```

**Tag Interest Score:**

```
tag_interest_score = sum_of_tag_matches / num_tags
```

### Real-World Examples

**User Interests:** Commented on "tech" (5x), "programming" (3x)

**Post A:** 1 hour old (0.04 days), 3 comments, tags: ["tech", "startup"]

```
recency_score = 1.0 / (1.0 + 0.04) = 0.96
engagement_score = min(3 / 10.0, 1.0) = 0.30
tag_interest_score = 0.5 (50% tag match)

final_score = (0.96 * 0.4) + (0.30 * 0.3) + (0.5 * 0.3)
            = 0.384 + 0.09 + 0.15
            = 0.624
```

**Post B:** 5 hours old (0.21 days), 8 comments, tags: ["random", "news"]

```
recency_score = 1.0 / (1.0 + 0.21) = 0.83
engagement_score = min(8 / 10.0, 1.0) = 0.80
tag_interest_score = 0.0 (no tag match)

final_score = (0.83 * 0.4) + (0.80 * 0.3) + (0.0 * 0.3)
            = 0.332 + 0.24 + 0.0
            = 0.572
```

**Post C:** 2 hours old (0.08 days), 0 comments, tags: ["programming", "tutorial"]

```
recency_score = 1.0 / (1.0 + 0.08) = 0.93
engagement_score = min(0 / 10.0, 1.0) = 0.0
tag_interest_score = 0.6 (60% tag match - related)

final_score = (0.93 * 0.4) + (0.0 * 0.3) + (0.6 * 0.3)
            = 0.372 + 0.0 + 0.18
            = 0.552
```

**Final Ranking:** A (0.624) > B (0.572) > C (0.552)

> Post A wins because it's fresh, has engagement, AND matches user interests!

---

## The Trade-off

The algorithm balances three factors:

- **Fresh content** (recent posts) - 40% weight
- **Popular content** (high engagement) - 30% weight
- **Personalized content** (matching user interests) - 30% weight

### Key Insights

- A fresh post matching user interests can beat an older popular post
- High engagement can overcome lack of personalization
- Very old posts need both high engagement AND tag matches to stay visible
- Posts with no tag matches rely entirely on recency and engagement

---

## Tuning the Algorithm

You can adjust the weights to change behavior:

### Adjust Component Weights

**Default (Balanced):**

```go
final_score = (recency * 0.4) + (engagement * 0.3) + (tag_interest * 0.3)
```

**More Weight on Recency (favor fresh content):**

```go
final_score = (recency * 0.5) + (engagement * 0.25) + (tag_interest * 0.25)
```

**More Weight on Engagement (favor popular content):**

```go
final_score = (recency * 0.3) + (engagement * 0.5) + (tag_interest * 0.2)
```

**More Weight on Personalization (favor user interests):**

```go
final_score = (recency * 0.3) + (engagement * 0.2) + (tag_interest * 0.5)
```

### Adjust Decay Rates

**Faster Recency Decay:**

```go
recency_score = 1.0 / (1.0 + age_in_hours)  // Decay by hour instead of day
```

**Higher Engagement Threshold:**

```go
engagement_score = min(comment_count / 20.0, 1.0)  // Need 20 comments for max score
```

### Adjust Related Tag Weight

**Stronger Related Tag Boost:**

```go
relatedTagScore = userInterest * 0.75  // 75% instead of 50%
```

**Weaker Related Tag Boost:**

```go
relatedTagScore = userInterest * 0.25  // 25% instead of 50%
```

---

## Implementation Guide

This section provides step-by-step instructions to implement the tag interest scoring feature.

### Current Implementation (Basic Feed)

Your current feed repository returns posts with basic scoring:

```go
// repositories/feed_repository.go
func (r *FeedRepository) GetUserFeed(ctx context.Context, userID int64) ([]*dtos.FeedPostResponseDTO, error) {
    query := `
        SELECT
            p.id,
            p.title,
            p.content,
            p.tags,
            u.id AS author_id,
            u.username,
            CONCAT(u.first_name, ' ', u.last_name) AS fullname,
            COUNT(c.id) AS comment_count,
            -- Basic scoring (old approach)
            (EXTRACT(EPOCH FROM NOW() - p.created_at) / 3600) * -1 + COUNT(c.id) * 10 AS score
        FROM posts p 
        INNER JOIN users u ON p.user_id = u.id
        LEFT JOIN comments c ON p.id = c.post_id 
        WHERE p.user_id IN (
            SELECT f.user_id
            FROM followers f
            WHERE f.follower_id = $1
        )
        GROUP BY p.id, u.id
        ORDER BY score DESC;
    `
    // ... scan and return
}
```

### Step 1: Add Method to Get User Tag Interests

Add a new method to `FeedRepository` to track which tags the user engages with:

```go
// repositories/feed_repository.go

func (r *FeedRepository) GetUserTagInterests(ctx context.Context, userID int64) (map[string]int, error) {
    ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
    defer cancel()
    
    query := `
        SELECT DISTINCT unnest(p.tags) as tag, COUNT(*) as engagement_count
        FROM comments c
        JOIN posts p ON c.post_id = p.id
        WHERE c.user_id = $1
        GROUP BY tag
        ORDER BY engagement_count DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, handleDBError(err, resourcePost)
    }
    defer rows.Close()
    
    interests := make(map[string]int)
    for rows.Next() {
        var tag string
        var count int
        if err := rows.Scan(&tag, &count); err != nil {
            return nil, handleDBError(err, resourcePost)
        }
        interests[tag] = count
    }
    
    return interests, rows.Err()
}
```

Update the interface in `domain/feed.go`:

```go
type FeedRepository interface {
    GetUserFeed(ctx context.Context, userID int64) ([]*dtos.FeedPostResponseDTO, error)
    GetUserTagInterests(ctx context.Context, userID int64) (map[string]int, error)
}
```

### Step 2: Update Feed Query to Return Normalized Scores

Modify the `GetUserFeed` query to return normalized scores instead of the old formula:

```go
// repositories/feed_repository.go

type FeedPostData struct {
    ID              int64
    Title           string
    Content         string
    Tags            []string
    CreatedAt       time.Time
    UpdatedAt       time.Time
    AuthorID        int64
    Username        string
    Fullname        string
    CommentCount    int
    RecencyScore    float64
    EngagementScore float64
}

func (r *FeedRepository) GetUserFeed(ctx context.Context, userID int64) ([]*FeedPostData, error) {
    ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
    defer cancel()
    
    query := `
        SELECT 
            p.id,
            p.title,
            p.content,
            p.tags,
            p.created_at,
            p.updated_at,
            u.id AS author_id,
            u.username,
            CONCAT(u.first_name, ' ', u.last_name) AS fullname,
            COUNT(DISTINCT c.id) as comment_count,
            
            -- Recency score (0-1, decays over time)
            1.0 / (1.0 + EXTRACT(EPOCH FROM (NOW() - p.created_at)) / 86400.0) as recency_score,
            
            -- Engagement score (0-1, normalized)
            LEAST(COUNT(DISTINCT c.id) / 10.0, 1.0) as engagement_score
            
        FROM posts p
        INNER JOIN followers f ON p.user_id = f.user_id
        INNER JOIN users u ON p.user_id = u.id
        LEFT JOIN comments c ON p.id = c.post_id
        WHERE f.follower_id = $1
        GROUP BY p.id, u.id, u.username, u.first_name, u.last_name
        ORDER BY p.created_at DESC
    `
    
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, handleDBError(err, resourcePost)
    }
    defer rows.Close()
    
    var feedData []*FeedPostData
    for rows.Next() {
        var data FeedPostData
        if err := rows.Scan(
            &data.ID,
            &data.Title,
            &data.Content,
            pq.Array(&data.Tags),
            &data.CreatedAt,
            &data.UpdatedAt,
            &data.AuthorID,
            &data.Username,
            &data.Fullname,
            &data.CommentCount,
            &data.RecencyScore,
            &data.EngagementScore,
        ); err != nil {
            return nil, handleDBError(err, resourcePost)
        }
        feedData = append(feedData, &data)
    }
    
    return feedData, rows.Err()
}
```

### Step 3: Define Related Tags Map

Create a map of related tags in your service layer:

```go
// services/feed_service.go

var relatedTags = map[string][]string{
    "store":       {"shop", "retail", "ecommerce", "shopping"},
    "tech":        {"technology", "programming", "software", "coding"},
    "food":        {"cooking", "recipe", "restaurant", "cuisine"},
    "travel":      {"vacation", "tourism", "adventure", "explore"},
    "fitness":     {"health", "workout", "exercise", "gym"},
    "business":    {"startup", "entrepreneur", "marketing", "sales"},
    "design":      {"ui", "ux", "graphics", "creative"},
    "education":   {"learning", "tutorial", "course", "teaching"},
    "gaming":      {"games", "esports", "streaming", "console"},
    "music":       {"audio", "concert", "band", "song"},
    "photography": {"photo", "camera", "portrait", "landscape"},
    "sports":      {"football", "basketball", "soccer", "athletics"},
    // Add more as needed based on your tags
}
```

### Step 4: Implement Tag Interest Scoring

Add helper functions to calculate tag interest scores:

```go
// services/feed_service.go

func (s *FeedService) calculateTagInterestScore(postTags []string, userInterests map[string]int) float64 {
    if len(postTags) == 0 || len(userInterests) == 0 {
        return 0.0
    }
    
    totalScore := 0.0
    maxEngagement := s.getMaxEngagement(userInterests)
    
    for _, postTag := range postTags {
        // Exact match: full weight
        if engagement, exists := userInterests[postTag]; exists {
            normalizedEngagement := float64(engagement) / float64(maxEngagement)
            totalScore += normalizedEngagement
        }
        
        // Related tags: partial weight (50%)
        for userTag, engagement := range userInterests {
            if s.areTagsRelated(postTag, userTag) {
                normalizedEngagement := float64(engagement) / float64(maxEngagement)
                totalScore += normalizedEngagement * 0.5
            }
        }
    }
    
    // Normalize by number of tags to avoid bias toward posts with many tags
    return totalScore / float64(len(postTags))
}

func (s *FeedService) areTagsRelated(tag1, tag2 string) bool {
    // Check if tag2 is in tag1's related tags
    if related, exists := relatedTags[tag1]; exists {
        for _, r := range related {
            if r == tag2 {
                return true
            }
        }
    }
    
    // Check if tag1 is in tag2's related tags
    if related, exists := relatedTags[tag2]; exists {
        for _, r := range related {
            if r == tag1 {
                return true
            }
        }
    }
    
    return false
}

func (s *FeedService) getMaxEngagement(interests map[string]int) int {
    max := 1 // Avoid division by zero
    for _, count := range interests {
        if count > max {
            max = count
        }
    }
    return max
}
```

### Step 5: Update Feed Service to Compute Final Scores

Modify the `GetUserFeed` method to combine all three scores:

```go
// services/feed_service.go

type scoredPost struct {
    data  *repositories.FeedPostData
    score float64
}

func (s *FeedService) GetUserFeed(ctx context.Context, userID int64) ([]*dtos.FeedPostResponseDTO, error) {
    // 1. Get user's tag interests
    tagInterests, err := s.feedRepository.GetUserTagInterests(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Get feed posts with base scores
    feedData, err := s.feedRepository.GetUserFeed(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Compute final scores with tag interest
    scored := make([]scoredPost, len(feedData))
    for i, post := range feedData {
        tagScore := s.calculateTagInterestScore(post.Tags, tagInterests)
        
        // Weighted combination: 40% recency, 30% engagement, 30% tag interest
        finalScore := (post.RecencyScore * 0.4) + 
                     (post.EngagementScore * 0.3) + 
                     (tagScore * 0.3)
        
        scored[i] = scoredPost{data: post, score: finalScore}
    }
    
    // 4. Sort by final score
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].score > scored[j].score
    })
    
    // 5. Convert to DTOs
    dtos := make([]*dtos.FeedPostResponseDTO, len(scored))
    for i, sp := range scored {
        dtos[i] = &dtos.FeedPostResponseDTO{
            ID:      sp.data.ID,
            Title:   sp.data.Title,
            Content: sp.data.Content,
            Tags:    sp.data.Tags,
            Author: dtos.AuthorInfoDTO{
                ID:       sp.data.AuthorID,
                Username: sp.data.Username,
                Fullname: sp.data.Fullname,
            },
            CreatedAt:    sp.data.CreatedAt,
            UpdatedAt:    sp.data.UpdatedAt,
            CommentCount: sp.data.CommentCount,
            Comments:     []dtos.CommentResponseDTO{}, // Empty for feed
            Score:        sp.score, // Include for debugging
        }
    }
    
    return dtos, nil
}
```

### Step 6: Testing

Test the implementation with different scenarios:

1. **User with no comment history** - Should fall back to recency + engagement only
2. **User who comments on "tech" posts** - Should see more tech-related posts
3. **User with diverse interests** - Should see balanced feed
4. **Posts with no tags** - Should still rank based on recency + engagement

### Summary

The implementation flow:

```
1. GetUserTagInterests() → map[string]int (tag → engagement count)
2. GetUserFeed() → []FeedPostData (with recency & engagement scores)
3. calculateTagInterestScore() → float64 (0-1 score per post)
4. Combine scores → final_score = (recency * 0.4) + (engagement * 0.3) + (tags * 0.3)
5. Sort by final_score → ranked feed
6. Convert to DTOs → return to handler
```

---

## Alternative Algorithms

### Reddit/Hacker News "Hot" Algorithm (Logarithmic Decay)

```sql
LOG(GREATEST(1, COUNT(c.id))) / POWER(
    (EXTRACT(EPOCH FROM NOW() - p.created_at) / 3600) + 2, 
    1.5
) as score
```

This creates **stronger decay** - older posts drop off faster even with high engagement.

---

## Future Enhancements

1. ✅ **Tag Preferences** - IMPLEMENTED: Boost posts with tags the user engages with
2. **Likes/Reactions** - Add when implemented
3. **User Affinity** - Boost posts from users you interact with frequently
4. **Time-of-Day Patterns** - Consider when user is most active
5. **Diversity** - Ensure variety in feed (not all from same user)
6. **Freshness Threshold** - Don't show posts older than X days
7. **Tag Co-occurrence** - Build dynamic related tags from data instead of manual mapping
8. **Negative Signals** - Downrank tags user has muted or hidden
9. **Collaborative Filtering** - "Users like you also engaged with..."

---

## Performance Considerations

### Required Indexes:
- `followers.follower_id` ✓
- `posts.created_at`
- `comments.post_id` ✓

### Caching Strategy:
- Cache feed for 1-5 minutes (Redis)
- Invalidate on new post from followed user
- Use cursor-based pagination for better performance

---

## Testing the Algorithm

### Scenarios to Test:

1. **Brand new post** (0 comments) vs **old popular post** (many comments)
2. **Multiple posts from same time** - should rank by engagement
3. **Empty feed** - user follows no one
4. **High-volume feed** - user follows many active users

### Metrics to Monitor:

- Average feed refresh rate
- Click-through rate on posts
- Time spent on feed
- Engagement rate (comments/likes per impression)
