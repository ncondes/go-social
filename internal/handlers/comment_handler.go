package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/dtos"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/metrics"
)

type CommentHandler struct {
	commentService domain.CommentServiceInterface
	validator      *Validator
	logger         logging.Logger
	metrics        *metrics.Metrics
}

func NewCommentHandler(
	commentService domain.CommentServiceInterface,
	validator *Validator,
	logger logging.Logger,
	metrics *metrics.Metrics,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		validator:      validator,
		logger:         logger,
		metrics:        metrics,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	user := getAuthenticatedUserFromContext(r.Context())
	postID := getPostIDFromContext(r.Context())

	var createCommentDTO *dtos.CreateCommentDTO

	if err := jsonDecode(w, r, &createCommentDTO); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	if err := h.validator.validateStruct(createCommentDTO); err != nil {
		respondWithErrors(w, http.StatusBadRequest, err, h.logger)
		return
	}

	comment := domain.Comment{
		Content: createCommentDTO.Content,
		UserID:  user.ID,
		PostID:  postID,
	}

	if err := h.commentService.CreateComment(r.Context(), &comment); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	h.metrics.CommentsCreated.Add(1)

	respondWithData(w, http.StatusCreated, comment, h.logger)
}

func (h *CommentHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	postID := getPostIDFromContext(r.Context())

	comments, err := h.commentService.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	responseDTOs := make([]*dtos.CommentResponseDTO, len(comments))
	for i, comment := range comments {
		responseDTOs[i] = new(dtos.CommentResponseDTO).FromDomain(comment)
	}

	respondWithData(w, http.StatusOK, responseDTOs, h.logger)
}
