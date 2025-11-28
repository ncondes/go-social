package services

import (
	"context"
	"sort"
	"sync"

	"github.com/ncondes/go/social/internal/domain"
)

type FeedService struct {
	feedRepository domain.FeedRepositoryInterface
}

func NewFeedService(feedRepository domain.FeedRepositoryInterface) *FeedService {
	return &FeedService{feedRepository: feedRepository}
}

func (s *FeedService) GetUserFeed(ctx context.Context, userID int64, options *domain.FeedQueryOptions) ([]*domain.FeedPost, error) {
	tagInterests, feed, err := s.fetchFeedData(ctx, userID, options)
	if err != nil {
		return nil, err
	}

	s.rankFeed(feed, tagInterests)

	return feed, nil
}

func (s *FeedService) fetchFeedData(ctx context.Context, userID int64, options *domain.FeedQueryOptions) (map[string]int, []*domain.FeedPost, error) {
	var (
		tagInterests map[string]int
		feed         []*domain.FeedPost
		tagErr       error
		feedErr      error
		wg           sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		tagInterests, tagErr = s.feedRepository.GetUserTagInterests(ctx, userID)
	}()

	go func() {
		defer wg.Done()
		feed, feedErr = s.feedRepository.GetUserFeed(ctx, userID, options)
	}()

	wg.Wait()

	if tagErr != nil {
		return nil, nil, tagErr
	}

	if feedErr != nil {
		return nil, nil, feedErr
	}

	return tagInterests, feed, nil
}

func (s *FeedService) rankFeed(feed []*domain.FeedPost, tagInterests map[string]int) {
	if len(feed) == 0 {
		return
	}

	maxEngagement := s.getMaxEngagement(tagInterests)

	for _, post := range feed {
		s.calculatePostScore(post, tagInterests, maxEngagement)
	}

	sort.Slice(feed, func(i, j int) bool {
		return feed[i].TotalScore > feed[j].TotalScore
	})
}

func (s *FeedService) getMaxEngagement(interests map[string]int) int {
	max := 1

	for _, count := range interests {
		if count > max {
			max = count
		}
	}

	return max
}

func (s *FeedService) calculatePostScore(post *domain.FeedPost, tagInterests map[string]int, maxEngagement int) {
	post.TagScore = s.calculateTagScore(post.Post.Tags, tagInterests, maxEngagement)
	post.TotalScore = s.calculateFinalScore(post.RecencyScore, post.EngagementScore, post.TagScore)
}

func (s *FeedService) calculateTagScore(postTags []string, userInterests map[string]int, maxEngagement int) float64 {
	if len(postTags) == 0 || len(userInterests) == 0 {
		return 0.0
	}

	totalScore := 0.0

	for _, postTag := range postTags {
		if engagement, exists := userInterests[postTag]; exists {
			totalScore += float64(engagement) / float64(maxEngagement)
		}
	}

	return totalScore / float64(len(postTags))
}

func (s *FeedService) calculateFinalScore(recency, engagement, tagScore float64) float64 {
	const (
		recencyWeight    = 0.4
		engagementWeight = 0.3
		tagWeight        = 0.3
	)

	return (recency * recencyWeight) + (engagement * engagementWeight) + (tagScore * tagWeight)
}
