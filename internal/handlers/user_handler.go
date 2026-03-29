package handlers

import (
	"net/http"

	"github.com/ncondes/go/social/internal/domain"
)

type UserHandler struct {
	userService domain.UserServiceInterface
}

func NewUserHandler(userService domain.UserServiceInterface) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := h.userService.CreateUser(r.Context(), &domain.User{
		FirstName: "Robertico",
	})
	if err != nil {
		handleError(w, r, err)
		return
	}

	respondWithData(w, http.StatusCreated, nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	respondWithData(w, http.StatusOK, user)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	// TODO: get follower ID from auth middleware in the future
	followerID := int64(1)

	if err := h.userService.FollowUser(r.Context(), userID, followerID); err != nil {
		handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	// TODO: get follower ID from auth middleware in the future
	followerID := int64(1)

	if err := h.userService.UnfollowUser(r.Context(), userID, followerID); err != nil {
		handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
