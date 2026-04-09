package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/domain"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/metrics"
	"github.com/ncondes/go/social/internal/ratelimit"
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

// Global rate limit (all requests)
func RateLimitMiddleware(rl ratelimit.RateLimiter, logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl == nil {
				next.ServeHTTP(w, r)
				return
			}

			if !checkRateLimit(r.Context(), w, rl, "global", logger) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Per-IP rate limit (separate limit per IP address)
func RateLimitByIPMiddleware(rl ratelimit.RateLimiter, logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl == nil {
				next.ServeHTTP(w, r)
				return
			}

			ip := getClientIP(r)
			key := "ip:" + ip

			if !checkRateLimit(r.Context(), w, rl, key, logger) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Per-user rate limit (for authenticated routes)
func RateLimitByAuthenticatedUserMiddleware(rl ratelimit.RateLimiter, logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if rl == nil {
				next.ServeHTTP(w, r)
				return
			}

			user := getAuthenticatedUserFromContext(r.Context())
			key := "user:" + strconv.FormatInt(user.ID, 10)

			if !checkRateLimit(r.Context(), w, rl, key, logger) {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	// Chi's middleware.RealIP has already processed X-Forwarded-For and X-Real-IP headers
	// and set r.RemoteAddr to the real client IP. We just need to strip the port.
	ip := r.RemoteAddr
	if host, _, found := strings.Cut(ip, ":"); found {
		return host
	}
	return ip
}

func checkRateLimit(
	ctx context.Context,
	w http.ResponseWriter,
	rl ratelimit.RateLimiter,
	key string,
	logger logging.Logger,
) bool {
	allowed, info, err := rl.Allow(ctx, key)

	// Fail closed - block request if rate limiter fails
	if err != nil {
		logger.Errorw("rate limit check failed",
			"error", err,
			"key", key,
		)
		respondWithError(w, http.StatusServiceUnavailable, "service temporarily unavailable", logger)
		return false
	}

	// Set rate limit headers
	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(info.Limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(info.Remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(info.Reset.Unix(), 10))

	if !allowed {
		w.Header().Set("Retry-After", strconv.FormatInt(info.Reset.Unix(), 10))
		logger.Warnw("rate limit exceeded",
			"key", key,
			"remaining", info.Remaining,
			"reset", info.Reset,
		)
		respondWithError(w, http.StatusTooManyRequests, "rate limit exceeded", logger)
		return false
	}

	return true
}

// MetricsMiddleware tracks HTTP request metrics for monitoring and observability.
// It measures request counts, errors, in-flight requests, and response times.
//
// How it works:
// 1. Records the start time when a request arrives
// 2. Increments total requests counter and in-flight requests counter
// 3. Wraps the ResponseWriter to capture the HTTP status code
// 4. Passes the request to the next handler in the chain
// 5. After the request completes:
//   - Decrements in-flight requests counter
//   - Records response time for this endpoint
//   - Increments error counter if status code >= 400
func MetricsMiddleware(m *metrics.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Record when the request started
			start := time.Now()

			// Track total requests received
			m.TotalRequests.Add(1)

			// Track how many requests are currently being processed
			m.RequestsInFlight.Add(1)
			// Ensure we decrement when this request finishes (even if it panics)
			defer m.RequestsInFlight.Add(-1)

			// Wrap the ResponseWriter so we can capture the status code
			// Default to 200 OK if WriteHeader is never called
			wrapped := &metricsResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Pass the request to the next handler
			next.ServeHTTP(wrapped, r)

			// Calculate how long the request took (in milliseconds)
			duration := time.Since(start).Milliseconds()
			// Record response time for this specific endpoint path
			m.ResponseTimes.Add(r.URL.Path, duration)

			// Count errors (4xx client errors and 5xx server errors)
			if wrapped.statusCode >= 400 {
				m.TotalErrors.Add(1)
			}
		})
	}
}

// metricsResponseWriter wraps http.ResponseWriter to capture the HTTP status code.
// This is necessary because the standard ResponseWriter doesn't expose the status code
// after it's been written.
type metricsResponseWriter struct {
	http.ResponseWriter     // Embed the original ResponseWriter
	statusCode          int // Store the status code when WriteHeader is called
}

// WriteHeader intercepts the status code before passing it to the underlying ResponseWriter.
// This method is automatically called by http.Handler when writing the response.
func (rw *metricsResponseWriter) WriteHeader(code int) {
	rw.statusCode = code                // Capture the status code for metrics
	rw.ResponseWriter.WriteHeader(code) // Pass it through to the real ResponseWriter
}
