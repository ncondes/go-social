package main

import (
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ncondes/go/social/internal/config"
	"github.com/ncondes/go/social/internal/db"
	dbPkg "github.com/ncondes/go/social/internal/db"
	"github.com/ncondes/go/social/internal/repositories"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// TODO: create seed load config function
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

	dbPkg.Flush(db)
	dbPkg.Seed(repositories)
}
