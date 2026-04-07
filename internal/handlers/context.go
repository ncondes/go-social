package handlers

import (
	"context"

	"github.com/ncondes/go/social/internal/domain"
)

type contextKey string

const (
	postIDContextKey            contextKey = "postID"
	userIDContextKey            contextKey = "userID"
	authenticatedUserContextKey contextKey = "authenticatedUser"
)

func getPostIDFromContext(ctx context.Context) int64 {
	postID, _ := ctx.Value(postIDContextKey).(int64)
	return postID
}

func getUserIDFromContext(ctx context.Context) int64 {
	userID, _ := ctx.Value(userIDContextKey).(int64)
	return userID
}

func getAuthenticatedUserFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(authenticatedUserContextKey).(*domain.User)
	return user
}
