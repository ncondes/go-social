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

// CreateComment godoc
//
//	@Summary		Create a comment on a post
//	@Description	Add a comment to a post by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64					true	"Post ID"
//	@Param			body	body		dtos.CreateCommentDTO	true	"Comment data"
//	@Success		201		{object}	dtos.CommentResponseDTO
//	@Failure		400		{object}	dtos.ErrorResponseDTO	"Bad request"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Post not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID}/comments [post]
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

// GetCommentsByPostID godoc
//
//	@Summary		Get comments for a post
//	@Description	Get all comments for a post by ID
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64	true	"Post ID"
//	@Success		200		{array}		dtos.CommentResponseDTO
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"Post not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/posts/{postID}/comments [get]
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
