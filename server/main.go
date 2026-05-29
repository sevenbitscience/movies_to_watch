package main

import (
	"log"

	"github.com/sevenbitscience/movies_to_watch/server/internal/database"
	"github.com/sevenbitscience/movies_to_watch/server/internal/web"
	"github.com/sevenbitscience/movies_to_watch/server/internal/tmdb"
	"github.com/sevenbitscience/movies_to_watch/server/internal/config"

)

func main() {
	// Load the configuration
	err := config.InitConfig(); if err != nil {
		log.Println("Failed to load config: ", err)
	}
	// Pass the TMDB API key over to tmdb.go
	tmdb.SetApiKey(*config.Config.TMBDkey)
	log.Println("Got TMDB API key")

	// Set up the DB
	database.InitDB(*config.Config.DatabasePath)
	log.Println("Connected to database")

	// Set up the REST API
	server := web.SetupServer(*config.Config.ServerURL)

	log.Printf("Server running on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
