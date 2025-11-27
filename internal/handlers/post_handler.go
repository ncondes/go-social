package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/services"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var createPostDTO *dtos.CreatePostDTO

	if err := jsonDecode(w, r, &createPostDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateStruct(createPostDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err)
		return
	}

	post := domain.Post{
		Title:   createPostDTO.Title,
		Content: createPostDTO.Content,
		Tags:    createPostDTO.Tags,
		UserID:  1, // TODO: get user ID from auth middleware in the future
	}

	if err := h.postService.CreatePost(r.Context(), &post); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusCreated, post); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postID := GetPostIDFromContext(r.Context())

	post, err := h.postService.GetPost(r.Context(), postID)
	if err != nil {
		switch err {
		case domain.ErrPostNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusOK, post); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postID := GetPostIDFromContext(r.Context())

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

	if err := validateStruct(updatePostDTO); err != nil {
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
		post.Tags = *updatePostDTO.Tags
	}

	if err := h.postService.UpdatePost(r.Context(), &post); err != nil {
		switch err {
		case domain.ErrPostNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusOK, &post); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := GetPostIDFromContext(r.Context())

	if err := h.postService.DeletePost(r.Context(), postID); err != nil {
		switch err {
		case domain.ErrPostNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
