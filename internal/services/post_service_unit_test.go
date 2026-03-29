package services

import (
	"context"
	"testing"
	"time"

	"github.com/ncondes/go/social/internal/domain"
)

func TestPostService_CreatePost(t *testing.T) {
	t.Run("should create post successfully", func(t *testing.T) {
		expectedID := int64(1)
		expectedCreatedAt := time.Now()
		expectedUpdatedAt := time.Now()

		mockPostRepository := &mockPostRepository{
			createFunc: func(ctx context.Context, post *domain.Post) error {
				// Simulate database creation setting these fields
				post.ID = expectedID
				post.CreatedAt = expectedCreatedAt
				post.UpdatedAt = expectedUpdatedAt
				return nil
			},
		}

		postService := NewPostService(mockPostRepository)

		post := &domain.Post{
			UserID:  1,
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		if err := postService.CreatePost(context.Background(), post); err != nil {
			t.Fatalf("failed to create post: %v", err)
		}

		if post.ID != expectedID {
			t.Errorf("expected post ID to be %d, got %d", expectedID, post.ID)
		}

		if !post.CreatedAt.Equal(expectedCreatedAt) {
			t.Errorf("expected post CreatedAt to be %v, got %v", expectedCreatedAt, post.CreatedAt)
		}

		if !post.UpdatedAt.Equal(expectedUpdatedAt) {
			t.Errorf("expected post UpdatedAt to be %v, got %v", expectedUpdatedAt, post.UpdatedAt)
		}

		if mockPostRepository.createCallCount != 1 {
			t.Errorf("expected create to be called once, got %d", mockPostRepository.createCallCount)
		}
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockPostRepository := &mockPostRepository{
			createFunc: func(ctx context.Context, post *domain.Post) error {
				return domain.ErrUserNotFound
			},
		}

		postService := NewPostService(mockPostRepository)

		post := &domain.Post{
			UserID:  1,
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		err := postService.CreatePost(context.Background(), post)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}

		if mockPostRepository.createCallCount != 1 {
			t.Errorf("expected create to be called once, got %d", mockPostRepository.createCallCount)
		}
	})

}

type mockPostRepository struct {
	createFunc       func(ctx context.Context, post *domain.Post) error
	getByIDFunc      func(ctx context.Context, postID int64) (*domain.PostWithDetails, error)
	updateFunc       func(ctx context.Context, post *domain.Post) error
	deleteFunc       func(ctx context.Context, postID int64) error
	createCallCount  int
	getByIDCallCount int
	updateCallCount  int
	deleteCallCount  int
}

func (m *mockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	m.createCallCount++
	if m.createFunc != nil {
		return m.createFunc(ctx, post)
	}
	return nil
}

func (m *mockPostRepository) GetByID(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
	m.getByIDCallCount++
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, postID)
	}
	return nil, nil
}

func (m *mockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	m.updateCallCount++
	if m.updateFunc != nil {
		return m.updateFunc(ctx, post)
	}
	return nil
}

func (m *mockPostRepository) Delete(ctx context.Context, postID int64) error {
	m.deleteCallCount++
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, postID)
	}
	return nil
}

var _ domain.PostRepositoryInterface = (*mockPostRepository)(nil)
