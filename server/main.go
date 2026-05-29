package main

import (
	"log"
	"os"

	"github.com/sevenbitscience/movies_to_watch/server/internal/database"
	"github.com/sevenbitscience/movies_to_watch/server/internal/rest"
	"github.com/sevenbitscience/movies_to_watch/server/internal/tmdb"

	"github.com/joho/godotenv" // Fetch info from .env file
)

func main() {
	// Fetch arguments from the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Pass the TMDB API key over to tmdb.go
	tmdb.SetApiKey(os.Getenv("TMDB_API_KEY"))
	log.Println("Got TMDB API key")

	// Set up the DB
	database.InitDB(os.Getenv("MOVIES_DATABASE_PATH"))
	log.Println("Connected to database")

	// Set up the REST API
	server := rest.SetupServer(os.Getenv("MOVIES_SERVER_URL"))

	log.Printf("Server running on %s", server.Addr)
	server.ListenAndServe()
}
