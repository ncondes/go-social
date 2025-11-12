package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/handlers"
)

type application struct {
	config   *config.Config
	handlers *handlers.Handlers
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID) // Add a request ID to the context
	r.Use(middleware.RealIP)    // Get the real IP address of the client
	r.Use(middleware.Logger)    // Log the request
	r.Use(middleware.Recoverer) // Recover from panics without crashing the server

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.handlers.Health.Check)
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	log.Printf("Starting server on %s\n", app.config.Addr)

	return srv.ListenAndServe()
}
