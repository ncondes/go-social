package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPostHandler_CreatePost(t *testing.T) {
	t.Parallel()

	t.Run("returns 201 when post is created successfully", func(t *testing.T) {
		t.Parallel()
		expectedID := int64(1)
		expectedCreatedAt := time.Now()
		expectedUpdatedAt := time.Now()

		mockPostService := &mockPostService{
			createPostFunc: func(ctx context.Context, post *domain.Post) error {
				post.ID = expectedID
				post.CreatedAt = expectedCreatedAt
				post.UpdatedAt = expectedUpdatedAt
				return nil
			},
		}

		validator := NewValidator()

		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		createPostDTO := dtos.CreatePostDTO{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", createPostDTO)
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Data domain.Post `json:"data"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, expectedID, response.Data.ID)
		assert.True(t, response.Data.CreatedAt.Equal(expectedCreatedAt))
		assert.True(t, response.Data.UpdatedAt.Equal(expectedUpdatedAt))
	})

	t.Run("returns 201 when post has empty tags and is created successfully", func(t *testing.T) {
		t.Parallel()
		expectedID := int64(1)
		expectedCreatedAt := time.Now()
		expectedUpdatedAt := time.Now()

		mockPostService := &mockPostService{
			createPostFunc: func(ctx context.Context, post *domain.Post) error {
				post.ID = expectedID
				post.CreatedAt = expectedCreatedAt
				post.UpdatedAt = expectedUpdatedAt
				return nil
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		createPostDTO := dtos.CreatePostDTO{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{},
		}

		req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", createPostDTO)
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response struct {
			Data domain.Post `json:"data"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, expectedID, response.Data.ID)
		assert.True(t, response.Data.CreatedAt.Equal(expectedCreatedAt))
		assert.True(t, response.Data.UpdatedAt.Equal(expectedUpdatedAt))
	})

	t.Run("returns 400 when JSON is invalid", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 400 when validation fails", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		testCases := []struct {
			name          string
			payload       dtos.CreatePostDTO
			expectedError string
		}{
			{
				name: "empty title",
				payload: dtos.CreatePostDTO{
					Title:   "",
					Content: "test content",
					Tags:    []string{"test"},
				},
				expectedError: "title is required",
			},
			{
				name: "empty content",
				payload: dtos.CreatePostDTO{
					Title:   "test title",
					Content: "",
					Tags:    []string{"test"},
				},
				expectedError: "content is required",
			},
			{
				name: "title too long",
				payload: dtos.CreatePostDTO{
					Title:   strings.Repeat("a", 101),
					Content: "test content",
					Tags:    []string{"test"},
				},
				expectedError: "title must be at most 100 characters",
			},
			{
				name: "content too long",
				payload: dtos.CreatePostDTO{
					Title:   "test title",
					Content: strings.Repeat("a", 1001),
					Tags:    []string{"test"},
				},
				expectedError: "content must be at most 1000 characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", tc.payload)
				w := httptest.NewRecorder()

				postHandler.CreatePost(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)

				var response struct {
					Errors []string `json:"errors"`
				}

				testutils.ParseJSONResponse(t, w, &response)

				assert.Contains(t, response.Errors, tc.expectedError)
			})
		}
	})

	t.Run("deduplicates tags when creating post", func(t *testing.T) {
		t.Parallel()

		var capturedPost *domain.Post
		mockPostService := &mockPostService{
			createPostFunc: func(ctx context.Context, post *domain.Post) error {
				capturedPost = post
				post.ID = 1
				return nil
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		createPostDTO := dtos.CreatePostDTO{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"go", "golang", "go", "backend", "golang"},
		}

		req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", createPostDTO)
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, []string{"go", "golang", "backend"}, capturedPost.Tags)
	})

	t.Run("returns 404 when user is not found", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			createPostFunc: func(ctx context.Context, post *domain.Post) error {
				return domain.ErrUserNotFound
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		createPostDTO := dtos.CreatePostDTO{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", createPostDTO)
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, domain.ErrUserNotFound.Error(), response.Error)
	})

	t.Run("returns 500 when internal error occurs", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			createPostFunc: func(ctx context.Context, post *domain.Post) error {
				return errors.New("database connection error")
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		createPostDTO := dtos.CreatePostDTO{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
		}

		req := testutils.MakeJSONRequest(t, http.MethodPost, "/posts", createPostDTO)
		w := httptest.NewRecorder()

		postHandler.CreatePost(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, "internal server error", response.Error)
	})
}

func TestPostHandler_GetPost(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with post details when post exists", func(t *testing.T) {
		t.Parallel()

		post := domain.Post{
			Title:   "Test Post",
			Content: "This is a test content",
			Tags:    []string{"test"},
			UserID:  1,
		}
		author := domain.User{
			ID:        1,
			Username:  "testuser",
			FirstName: "Test",
			LastName:  "User",
		}

		mockPostService := &mockPostService{
			getPostFunc: func(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
				return &domain.PostWithDetails{
					Post:         post,
					Author:       author,
					CommentCount: 2,
				}, nil
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodGet, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.GetPost(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data *dtos.PostResponseDTO `json:"data"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.NotNil(t, response.Data)
		assert.Equal(t, post.Title, response.Data.Title)
		assert.Equal(t, post.Content, response.Data.Content)
		assert.Equal(t, post.Tags, response.Data.Tags)
		assert.Equal(t, author.ID, response.Data.Author.ID)
		assert.Equal(t, author.Username, response.Data.Author.Username)
		assert.Equal(t, author.FullName(), response.Data.Author.Fullname)
		assert.Equal(t, 2, response.Data.CommentCount)
	})

	t.Run("returns 404 when post is not found", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			getPostFunc: func(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
				return nil, domain.ErrPostNotFound
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodGet, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.GetPost(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, domain.ErrPostNotFound.Error(), response.Error)
	})

	t.Run("returns 500 when internal error occurs", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			getPostFunc: func(ctx context.Context, postID int64) (*domain.PostWithDetails, error) {
				return nil, errors.New("database connection error")
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodGet, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.GetPost(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, "internal server error", response.Error)
	})
}

func TestPostHandler_UpdatePost(t *testing.T) {
	t.Parallel()

	t.Run("returns 200 with updated post when update succeeds", func(t *testing.T) {
		t.Parallel()

		title := "Updated Title"
		content := "Updated Content"
		tags := []string{"updated", "tags"}
		version := int64(1)

		mockPostService := &mockPostService{
			updatePostFunc: func(ctx context.Context, post *domain.Post) error {
				post.Title = title
				post.Content = content
				post.Tags = tags
				post.Version = version + 1
				return nil
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		updatePostDTO := dtos.UpdatePostDTO{
			Title:   &title,
			Content: &content,
			Tags:    &tags,
			Version: &version,
		}

		req := testutils.MakeJSONRequest(t, http.MethodPatch, "/posts/1", updatePostDTO)
		w := httptest.NewRecorder()

		postHandler.UpdatePost(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Data *domain.Post `json:"data"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.NotNil(t, response.Data)
		assert.Equal(t, title, response.Data.Title)
		assert.Equal(t, content, response.Data.Content)
		assert.Equal(t, tags, response.Data.Tags)
		assert.Equal(t, version+1, response.Data.Version)
	})

	t.Run("returns 400 when JSON is invalid", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := httptest.NewRequest(http.MethodPatch, "/posts/1", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		postHandler.UpdatePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 400 when no fields are provided to update", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		updatePostDTO := dtos.UpdatePostDTO{}

		req := testutils.MakeJSONRequest(t, http.MethodPatch, "/posts/1", updatePostDTO)
		w := httptest.NewRecorder()

		postHandler.UpdatePost(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, "no fields to update", response.Error)
	})

	t.Run("returns 400 when validation fails", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		testCases := []struct {
			name          string
			payload       dtos.UpdatePostDTO
			expectedError string
		}{
			{
				name: "empty title",
				payload: dtos.UpdatePostDTO{
					Title: func() *string {
						str := ""
						return &str
					}(),
					Version: func() *int64 {
						v := int64(1)
						return &v
					}(),
				},
				expectedError: "title must be at least 1 characters",
			},
			{
				name: "empty content",
				payload: dtos.UpdatePostDTO{
					Content: func() *string {
						str := ""
						return &str
					}(),
					Version: func() *int64 {
						v := int64(1)
						return &v
					}(),
				},
				expectedError: "content must be at least 1 characters",
			},
			{
				name: "title too long",
				payload: dtos.UpdatePostDTO{
					Title: func() *string {
						str := strings.Repeat("a", 101)
						return &str
					}(),
					Version: func() *int64 {
						v := int64(1)
						return &v
					}(),
				},
				expectedError: "title must be at most 100 characters",
			},
			{
				name: "content too long",
				payload: dtos.UpdatePostDTO{
					Content: func() *string {
						str := strings.Repeat("a", 1001)
						return &str
					}(),
					Version: func() *int64 {
						v := int64(1)
						return &v
					}(),
				},
				expectedError: "content must be at most 1000 characters",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				req := testutils.MakeJSONRequest(t, http.MethodPatch, "/posts/1", tc.payload)
				w := httptest.NewRecorder()

				postHandler.UpdatePost(w, req)

				assert.Equal(t, http.StatusBadRequest, w.Code)

				var response struct {
					Errors []string `json:"errors"`
				}

				testutils.ParseJSONResponse(t, w, &response)

				assert.Contains(t, response.Errors, tc.expectedError)
			})
		}
	})

	t.Run("returns 404 when post is not found", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			updatePostFunc: func(ctx context.Context, post *domain.Post) error {
				return domain.ErrPostNotFound
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodPatch, "/posts/1", dtos.UpdatePostDTO{
			Title: func() *string {
				str := "test"
				return &str
			}(),
			Version: func() *int64 {
				v := int64(1)
				return &v
			}(),
		})
		w := httptest.NewRecorder()

		postHandler.UpdatePost(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, domain.ErrPostNotFound.Error(), response.Error)
	})

	t.Run("returns 500 when internal error occurs", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			updatePostFunc: func(ctx context.Context, post *domain.Post) error {
				return errors.New("database connection error")
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		title := "test"
		version := int64(1)
		req := testutils.MakeJSONRequest(t, http.MethodPatch, "/posts/1", dtos.UpdatePostDTO{
			Title:   &title,
			Version: &version,
		})
		w := httptest.NewRecorder()

		postHandler.UpdatePost(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, "internal server error", response.Error)
	})
}

func TestPostHandler_DeletePost(t *testing.T) {
	t.Parallel()

	t.Run("returns 204 when post is deleted successfully", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			deletePostFunc: func(ctx context.Context, postID int64) error {
				return nil
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodDelete, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.DeletePost(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("returns 404 when post is not found", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			deletePostFunc: func(ctx context.Context, postID int64) error {
				return domain.ErrPostNotFound
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodDelete, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.DeletePost(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, domain.ErrPostNotFound.Error(), response.Error)
	})

	t.Run("returns 500 when internal error occurs", func(t *testing.T) {
		t.Parallel()

		mockPostService := &mockPostService{
			deletePostFunc: func(ctx context.Context, postID int64) error {
				return errors.New("database connection error")
			},
		}

		validator := NewValidator()
		postHandler := NewPostHandler(mockPostService, validator, testLogger)

		req := testutils.MakeJSONRequest(t, http.MethodDelete, "/posts/1", nil)
		w := httptest.NewRecorder()

		postHandler.DeletePost(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response struct {
			Error string `json:"error"`
		}

		testutils.ParseJSONResponse(t, w, &response)

		assert.Equal(t, "internal server error", response.Error)
	})
}
