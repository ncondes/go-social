package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/logging"
)

type PostHandler struct {
	postService domain.PostServiceInterface
	validator   *Validator
	logger      logging.Logger
	authorizer  *auth.Authorizer
}

func NewPostHandler(
	postService domain.PostServiceInterface,
	validator *Validator,
	logger logging.Logger,
	authorizer *auth.Authorizer,
) *PostHandler {
	return &PostHandler{
		postService: postService,
		validator:   validator,
		logger:      logger,
		authorizer:  authorizer,
	}
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content, and optional tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		dtos.CreatePostDTO	true	"Post data"
//	@Success		201		{object}	domain.Post
//	@Failure		400		{object}	dtos.ErrorsResponseDTO	"Validation errors"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts [post]
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var createPostDTO *dtos.CreatePostDTO

	if err := jsonDecode(w, r, &createPostDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := h.validator.validateStruct(createPostDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err, h.logger)
		return
	}

	user := getAuthenticatedUserFromContext(r.Context())

	post := domain.Post{
		Title:   createPostDTO.Title,
		Content: createPostDTO.Content,
		Tags:    h.deduplicateTags(createPostDTO.Tags),
		UserID:  user.ID,
	}

	if err := h.postService.CreatePost(r.Context(), &post); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusCreated, post, h.logger)
}

// GetPost godoc
//
//	@Summary		Get a post
//	@Description	Get a post by ID with author details and comment count
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64	true	"Post ID"
//	@Success		200		{object}	dtos.PostResponseDTO
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Post not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID} [get]
func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	postWithDetails, err := h.postService.GetPost(r.Context(), postID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	response := new(dtos.PostResponseDTO).FromDomain(postWithDetails)

	respondWithData(w, http.StatusOK, response, h.logger)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update a post's title, content, or tags by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64				true	"Post ID"
//	@Param			post	body		dtos.UpdatePostDTO	true	"Fields to update"
//	@Success		200		{object}	domain.Post
//	@Failure		400		{object}	dtos.ErrorsResponseDTO	"Validation errors"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		403		{object}	dtos.ErrorResponseDTO	"Insufficient permissions"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Post not found"
//	@Failure		409		{object}	dtos.ErrorResponseDTO	"Post version conflict"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID} [patch]
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	user := getAuthenticatedUserFromContext(r.Context())
	postID := getPostIDFromContext(r.Context())

	var updatePostDTO *dtos.UpdatePostDTO

	if err := jsonDecode(w, r, &updatePostDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	// TODO: should we make this shared? as a util
	if updatePostDTO.Title == nil && updatePostDTO.Content == nil && updatePostDTO.Tags == nil {
		respondWithError(w, http.StatusBadRequest, "no fields to update", h.logger)
		return
	}

	if err := h.validator.validateStruct(updatePostDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err, h.logger)
		return
	}

	existingPost, err := h.postService.GetPost(r.Context(), postID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	if !h.authorizer.CanUpdatePost(user, &existingPost.Post) {
		respondWithError(w, http.StatusForbidden, "insufficient permissions", h.logger)
		return
	}

	post := domain.Post{
		ID:      postID,
		Version: *updatePostDTO.Version,
	}

	if updatePostDTO.Title != nil {
		post.Title = *updatePostDTO.Title
	}

	if updatePostDTO.Content != nil {
		post.Content = *updatePostDTO.Content
	}

	if updatePostDTO.Tags != nil {
		post.Tags = h.deduplicateTags(*updatePostDTO.Tags)
	}

	if err := h.postService.UpdatePost(r.Context(), &post); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusOK, &post, h.logger)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete a post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int64	true	"Post ID"
//	@Success		204		"No content"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		403		{object}	dtos.ErrorResponseDTO	"Insufficient permissions"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Post not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID} [delete]
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	user := getAuthenticatedUserFromContext(r.Context())
	postID := getPostIDFromContext(r.Context())

	existingPost, err := h.postService.GetPost(r.Context(), postID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	if !h.authorizer.CanDeletePost(user, &existingPost.Post) {
		respondWithError(w, http.StatusForbidden, "insufficient permissions", h.logger)
		return
	}

	if err := h.postService.DeletePost(r.Context(), postID); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PostHandler) deduplicateTags(tags []string) []string {
	if len(tags) == 0 {
		return tags
	}

	seen := make(map[string]bool, len(tags))
	result := make([]string, 0, len(tags))

	for _, tag := range tags {
		if !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}

	return result
}
