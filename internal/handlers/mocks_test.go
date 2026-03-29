package handlers

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type mockPostService struct {
	createPostFunc      func(ctx context.Context, post *domain.Post) error
	getPostFunc         func(ctx context.Context, postID int64) (*domain.PostWithDetails, error)
	updatePostFunc      func(ctx context.Context, post *domain.Post) error
	deletePostFunc      func(ctx context.Context, postID int64) error
	createPostCallCount int
	getPostCallCount    int
	updatePostCallCount int
	deletePostCallCount int
}

func (m *mockPostService) CreatePost(ctx context.Context, post *domain.Post) error {
	m.createPostCallCount++

	if m.createPostFunc != nil {
		return m.createPostFunc(ctx, post)
	}

	return nil
}

func (m *mockPostService) GetPost(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
	m.getPostCallCount++

	if m.getPostFunc != nil {
		return m.getPostFunc(ctx, postID)
	}

	return nil, nil
}

func (m *mockPostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	m.updatePostCallCount++

	if m.updatePostFunc != nil {
		return m.updatePostFunc(ctx, post)
	}

	return nil
}

func (m *mockPostService) DeletePost(ctx context.Context, postID int64) error {
	m.deletePostCallCount++

	if m.deletePostFunc != nil {
		return m.deletePostFunc(ctx, postID)
	}

	return nil
}

var _ domain.PostServiceInterface = (*mockPostService)(nil)
