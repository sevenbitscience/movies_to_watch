package main

import (
	"log"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Movie struct {
	ID 		int 	`json:"id"`
	Title	string	`json:"title"`
	TMDB_ID	int		`json:"tmdb_id"`
	Year	int		`json:"year"`
	Genres	string	`json:"genres"`
	Status	string	`json:"status"`
}

const movieSchema = `
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tmdb_id INTEGER NOT NULL UNIQUE,
    title TEXT NOT NULL,
    release_year INTEGER,
    genres TEXT,
    status TEXT DEFAULT 'pending'
);`

var db *sql.DB

/* Initiate the connection to the database
 * 
 * Creates a new database if there is not one at the provided path.
 * Takes the path to the sqlite database the parameter
 */
func initDB(path string) {
	var err error

	db, err = sql.Open("sqlite", path)
	if (err != nil) {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	
	_, err = db.Exec(movieSchema)
	if (err != nil) {
		log.Fatalf("Couldn't write schema to database: %v", err)
	}
}

// Get all the movies from the database
// TODO: Get n movies from the database or filter movies from DB
func getWatchlist() []Movie {
	log.Fatalln("Not implemented")
	return nil
}


// Add a movie to the database
func addMovie(movie *Movie) {}


// Set a movie as watched
func setWatchedStatus(id int, status string) {}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyTMDB(tmdb_id int) bool {
	return false
}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyID(id int) bool {
	return false
}
