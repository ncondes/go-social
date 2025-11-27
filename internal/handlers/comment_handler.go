package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/services"
)

type CommentHandler struct {
	commentService *services.CommentService
}

func NewCommentHandler(commentService *services.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := GetPostIDFromContext(r.Context())

	var createCommentDTO *dtos.CreateCommentDTO

	if err := jsonDecode(w, r, &createCommentDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateStruct(createCommentDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err)
		return
	}

	comment := domain.Comment{
		Content: createCommentDTO.Content,
		UserID:  1, // TODO: get user ID from auth middleware in the future
		PostID:  postID,
	}

	if err := h.commentService.CreateComment(r.Context(), &comment); err != nil {
		switch err {
		case domain.ErrPostNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusCreated, comment); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *CommentHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postID := GetPostIDFromContext(r.Context())

	comments, err := h.commentService.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		handleInternalServerError(w, r, err)
		return
	}

	if err := respondWithData(w, http.StatusOK, comments); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}
