package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ncondes/go/social/internal/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	app := &application{
		config: cfg,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
