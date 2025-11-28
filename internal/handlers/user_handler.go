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
		switch err {
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusCreated, nil); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	if err := respondWithData(w, http.StatusOK, user); err != nil {
		handleInternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	// TODO: get follower ID from auth middleware in the future
	followerID := int64(1)

	if err := h.userService.FollowUser(r.Context(), userID, followerID); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		case domain.ErrUserAlreadyFollowing:
			// Idempotency handling
			respondWithError(w, http.StatusNoContent, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	// TODO: get follower ID from auth middleware in the future
	followerID := int64(1)

	if err := h.userService.UnfollowUser(r.Context(), userID, followerID); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		default:
			handleInternalServerError(w, r, err)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
