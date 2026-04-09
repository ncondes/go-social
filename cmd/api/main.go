package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/ncondes/go/social/internal/auth"
	"github.com/ncondes/go/social/internal/cache"
	cfg "github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/db"
	"github.com/ncondes/go/social/internal/handlers"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/mailer"
	"github.com/ncondes/go/social/internal/repositories"
	"github.com/ncondes/go/social/internal/services"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "1.0.0"

//	@title	Gopher Social API
//	@version
//	@description	A RESTful API for a social network. Supports user profiles,
//	@description	posts with tags, comments, follower relationships, and a personalized feed
//	@description	scored by recency, engagement, and tag affinity.

//	@contact.name	Nicolas Conde
//	@contact.url	https://github.com/ncondes/go/social

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func main() {
	logger := logging.NewZapLogger(zap.Must(zap.NewProduction()).Sugar())
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatalw("Error loading .env file", "error", err)
	}

	config := cfg.Load()

	db, err := db.New(
		config.DB.Addr,
		config.DB.MaxOpenConns,
		config.DB.MaxIdleConns,
		config.DB.MaxIdleTime,
	)
	if err != nil {
		logger.Fatalw("Error connecting to database", "error", err)
	}

	logger.Infow("Connected to database successfully")
	defer db.Close()

	var rc *redis.Client

	if config.Redis.Enabled {
		rc = cache.NewRedisClient(
			config.Redis.Addr,
			config.Redis.Password,
			config.Redis.DB,
		)

		logger.Infow("Connected to Redis successfully")
	} else {
		rc = nil
	}

	cacheStorage := cache.NewCacheStorage(rc)

	rateLimiters := newRateLimiters(rc, config.RateLimit)

	if config.RateLimit.Enabled {
		logger.Infow("Rate limiters initialized",
			"global", fmt.Sprintf("%d/%s", config.RateLimit.Global.RequestsPerWindow, config.RateLimit.Global.Window),
			"strict_ip", fmt.Sprintf("%d/%s", config.RateLimit.StrictIP.RequestsPerWindow, config.RateLimit.StrictIP.Window),
			"read_ops", fmt.Sprintf("%d/%s", config.RateLimit.ReadOps.RequestsPerWindow, config.RateLimit.ReadOps.Window),
			"write_ops", fmt.Sprintf("%d/%s", config.RateLimit.WriteOps.RequestsPerWindow, config.RateLimit.WriteOps.Window),
		)
	} else {
		logger.Infow("Rate limiters disabled")
	}

	repositories := repositories.New(db, config)
	mailer := mailer.NewSendGridMailer(config.MailConfig.FromEmail, config.MailConfig.APIKey)
	authenticator := auth.NewJWTAuthenticator(
		config.Auth.JWT.Secret,
		config.Auth.JWT.Audience,
		config.Auth.JWT.Issuer,
		config.Auth.JWT.Duration,
	)
	services := services.New(
		repositories,
		config,
		mailer,
		logger,
		authenticator,
		cacheStorage,
	)
	validator := handlers.NewValidator()
	authorizer := auth.NewAuthorizer()
	handlers := handlers.New(config,
		services,
		validator,
		logger,
		authorizer,
	)

	app := &application{
		config:        config,
		handlers:      handlers,
		logger:        logger,
		services:      services,
		authenticator: authenticator,
		rateLimiters:  rateLimiters,
	}

	mux := app.mount()

	if err := app.run(mux); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Infow("Server stopped")
			return
		}

		logger.Fatalw("Server error", "error", err)
	}
}
