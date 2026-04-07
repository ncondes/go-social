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

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	respondWithData(w, http.StatusOK, user, h.logger)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	follower := getAuthenticatedUserFromContext(r.Context())

	if err := h.userService.FollowUser(r.Context(), userID, follower.ID); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	follower := getAuthenticatedUserFromContext(r.Context())

	if err := h.userService.UnfollowUser(r.Context(), userID, follower.ID); err != nil {
		handleError(w, r, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
