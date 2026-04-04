package main

import (
	"errors"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/db"
	"github.com/ncondes/go/social/internal/handlers"
	"github.com/ncondes/go/social/internal/logging"
	"github.com/ncondes/go/social/internal/repositories"
	"github.com/ncondes/go/social/internal/services"
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

	config := config.Load()

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

	repositories := repositories.New(db, config)
	services := services.New(repositories)
	validator := handlers.NewValidator()
	handlers := handlers.New(config, services, validator, logger)

	app := &application{
		config:   config,
		handlers: handlers,
		logger:   logger,
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
