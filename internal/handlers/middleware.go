package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
)

func PostIDMiddleware(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			postIDParam := chi.URLParam(r, "postID")
			postID, err := strconv.ParseInt(postIDParam, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, "invalid post ID", logger)
				return
			}

			ctx := context.WithValue(r.Context(), postIDContextKey, postID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDMiddleware(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDParam := chi.URLParam(r, "userID")
			userID, err := strconv.ParseInt(userIDParam, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid credentials", logger)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BasicAuthMiddleware(logger logging.Logger, config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header", logger)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header", logger)
				return
			}

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header", logger)
				return
			}

			username := config.Auth.Basic.Username
			password := config.Auth.Basic.Password

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header", logger)
				return
			}

			if credentials[0] != username || credentials[1] != password {
				respondWithError(w, http.StatusUnauthorized, "invalid credentials", logger)
				return
			}

			ctx := context.Background()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthTokenMiddleware(
	authenticator *auth.JWTAuthenticator,
	userService domain.UserServiceInterface,
	logger logging.Logger,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "missing authorization header", logger)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "invalid authorization header", logger)
				return
			}

			token := parts[1]
			// Validate the token and get claims
			claims, err := authenticator.ValidateToken(token)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid or expired token", logger)
				return
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid token claims", logger)
				return
			}

			logger.Infow("user ID from token", "userID", userID)

			user, err := userService.GetUser(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "invalid credentials", logger)
				return
			}

			ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
