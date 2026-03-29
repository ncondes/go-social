package repositories

import (
	"context"
	"testing"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/testutils"
)

func TestPostRepository_Create(t *testing.T) {
	testDB := testutils.NewTestDB(t)
	defer testutils.TeardownTestDB(t, testDB)

	postRepository := NewPostRepository(testDB)

	t.Run("should create post successfully", func(t *testing.T) {
		testutils.TruncateTables(t, testDB, "posts", "users")

		userID := testutils.CreateTestUser(t, testDB)

		post := &domain.Post{
			UserID:  userID,
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		if err := postRepository.Create(context.Background(), post); err != nil {
			t.Fatalf("failed to create post: %v", err)
		}

		if post.ID == 0 {
			t.Errorf("expected post ID to be set, got %d", post.ID)
		}

		if post.CreatedAt.IsZero() {
			t.Errorf("expected post CreatedAt to be set, got %v", post.CreatedAt)
		}

		if post.UpdatedAt.IsZero() {
			t.Errorf("expected post UpdatedAt to be set, got %v", post.UpdatedAt)
		}

		// Verify that the post was created in the database
		var count int
		if err := testDB.QueryRow("SELECT COUNT(*) FROM posts WHERE id = $1", post.ID).Scan(&count); err != nil {
			t.Fatalf("failed to query database: %v", err)
		}

		if count != 1 {
			t.Errorf("expected post to be created in the database, got %d", count)
		}
	})

	t.Run("should return error for non-existing user", func(t *testing.T) {
		var userID int64 = 99999999999

		post := &domain.Post{
			UserID:  userID,
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		err := postRepository.Create(context.Background(), post)

		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

var _ domain.PostRepositoryInterface = (*PostRepository)(nil)
