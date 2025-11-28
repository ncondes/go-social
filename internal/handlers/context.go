package handlers

import "context"

type contextKey string

const (
	postIDContextKey contextKey = "postID"
	userIDContextKey contextKey = "userID"
)

func getPostIDFromContext(ctx context.Context) int64 {
	postID, _ := ctx.Value(postIDContextKey).(int64)
	return postID
}

func getUserIDFromContext(ctx context.Context) int64 {
	userID, _ := ctx.Value(userIDContextKey).(int64)
	return userID
}
