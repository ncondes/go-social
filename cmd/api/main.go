package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/db"
	"github.com/ncondes/go/social/internal/handlers"
	"github.com/ncondes/go/social/internal/repositories"
	"github.com/ncondes/go/social/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.Load()

	db, err := db.New(
		config.DB.Addr,
		config.DB.MaxOpenConns,
		config.DB.MaxIdleConns,
		config.DB.MaxIdleTime,
	)
	if err != nil {
		log.Println(err)
		log.Fatal("Error connecting to database")
	}

	log.Println("Connected to database successfully")
	defer db.Close()

	repositories := repositories.New(db)
	services := services.New(repositories)
	validator := handlers.NewValidator()
	handlers := handlers.New(config, services, validator)

	app := &application{
		config:   config,
		handlers: handlers,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
