package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
)

type CommentHandler struct {
	commentService domain.CommentServiceInterface
	validator      *Validator
}

func NewCommentHandler(commentService domain.CommentServiceInterface, validator *Validator) *CommentHandler {
	return &CommentHandler{commentService: commentService, validator: validator}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	var createCommentDTO *dtos.CreateCommentDTO

	if err := jsonDecode(w, r, &createCommentDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validator.validateStruct(createCommentDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err)
		return
	}

	comment := domain.Comment{
		Content: createCommentDTO.Content,
		UserID:  1, // TODO: get user ID from auth middleware in the future
		PostID:  postID,
	}

	if err := h.commentService.CreateComment(r.Context(), &comment); err != nil {
		handleError(w, r, err)
		return
	}

	respondWithData(w, http.StatusCreated, comment)
}

func (h *CommentHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	comments, err := h.commentService.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	responseDTOs := make([]*dtos.CommentResponseDTO, len(comments))
	for i, comment := range comments {
		responseDTOs[i] = new(dtos.CommentResponseDTO).FromDomain(comment)
	}

	respondWithData(w, http.StatusOK, responseDTOs)
}
