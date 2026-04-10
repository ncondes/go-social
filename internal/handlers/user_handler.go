package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
)

type UserHandler struct {
	userService domain.UserServiceInterface
	logger      logging.Logger
}

func NewUserHandler(userService domain.UserServiceInterface, logger logging.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := h.userService.CreateUser(r.Context(), &domain.User{
		FirstName: "Robertico",
	})
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusCreated, nil, h.logger)
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64	true	"User ID"
//	@Success		200		{object}	domain.User
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"User not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusOK, user, h.logger)
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int64	true	"User ID to follow"
//	@Success		204		"No content"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"User not found"
//	@Failure		409		{object}	dtos.ErrorResponseDTO	"Already following"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID}/follow [post]
func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	follower := getAuthenticatedUserFromContext(r.Context())

	if err := h.userService.FollowUser(r.Context(), userID, follower.ID); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int64	true	"User ID to unfollow"
//	@Success		204		"No content"
//	@Failure		401		{object}	dtos.ErrorResponseDTO	"Unauthorized"
//	@Failure		404		{object}	dtos.ErrorResponseDTO	"User not found"
//	@Failure		500		{object}	dtos.ErrorResponseDTO	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID}/unfollow [delete]
func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	follower := getAuthenticatedUserFromContext(r.Context())

	if err := h.userService.UnfollowUser(r.Context(), userID, follower.ID); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
