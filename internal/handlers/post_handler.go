package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
)

type PostHandler struct {
	postService domain.PostServiceInterface
	validator   *Validator
}

func NewPostHandler(postService domain.PostServiceInterface, validator *Validator) *PostHandler {
	return &PostHandler{postService: postService, validator: validator}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var createPostDTO *dtos.CreatePostDTO

	if err := jsonDecode(w, r, &createPostDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.validateStruct(createPostDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err)
		return
	}

	post := domain.Post{
		Title:   createPostDTO.Title,
		Content: createPostDTO.Content,
		Tags:    h.deduplicateTags(createPostDTO.Tags),
		UserID:  1, // TODO: get user ID from auth middleware in the future
	}

	if err := h.postService.CreatePost(r.Context(), &post); err != nil {
		handleError(w, r, err)
		return
	}

	respondWithData(w, http.StatusCreated, post)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	postWithDetails, err := h.postService.GetPost(r.Context(), postID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response := new(dtos.PostResponseDTO).FromDomain(postWithDetails)

	respondWithData(w, http.StatusOK, response)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	var updatePostDTO *dtos.UpdatePostDTO

	if err := jsonDecode(w, r, &updatePostDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: should we make this shared? as a util
	if updatePostDTO.Title == nil && updatePostDTO.Content == nil && updatePostDTO.Tags == nil {
		respondWithError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	if err := h.validator.validateStruct(updatePostDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err)
		return
	}

	post := domain.Post{
		ID:        postID,
		UpdatedAt: *updatePostDTO.UpdatedAt,
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
		handleError(w, r, err)
		return
	}

	respondWithData(w, http.StatusOK, &post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	if err := h.postService.DeletePost(r.Context(), postID); err != nil {
		handleError(w, r, err)
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
