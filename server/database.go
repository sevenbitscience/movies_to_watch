package main

import (
	"database/sql"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type Movie struct {
	ID 		int 	`json:"id"`
	Title	string	`json:"title"`
	TMDB_ID	int		`json:"tmdb_id"`
	Year	int		`json:"release_year"`
	Genres	[]string	`json:"genre"`
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
// Pagination with LIMIT and OFFSET
func getWatchlist() ([]Movie, error) {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		log.Println("Couldn't get movies from database")
		return nil, err
	}
	defer rows.Close()

	var watchlist []Movie

	for rows.Next() {
		var m Movie
		err := rows.Scan(&m.ID, &m.TMDB_ID, &m.Title, &m.Year, &m.Genres, &m.Status)
		if err != nil {
			log.Println("Error reading movies from database!")
			return nil, err
		}
		watchlist = append(watchlist, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return watchlist, nil
}


// Add a movie to the database
func addMovie(movie *Movie) {
	genreString := strings.Join(movie.Genres, ", ")
	query := `
	INSERT INTO movies 
	(tmdb_id, title, release_year, genres)
	VALUES
	(?, ?, ?, ?)`
	//WHERE NOT EXISTS (
	//	SELECT 1 FROM movies WHERE tmdb_id = ?
	//)`
	_, err := db.Exec(query,
	movie.TMDB_ID,
	movie.Title,
	movie.Year,
	genreString,
	)

	if (err != nil) {
		log.Printf("Couldn't add movie %+v", movie)
	}
}


// Set a movie as watched
func setWatchedStatus(id int, status string) {
	query := "UPDATE movies SET status = ? WHERE id = ?"
	_, err := db.Exec(query, status, id)
	if (err != nil) {
		log.Printf("Failed to execute SQL", )
	}
}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyTMDB(tmdb_id int) bool {
	err := db.QueryRow("SELECT * FROM movies WHERE tmdb_id = ?", tmdb_id)
	return err == nil 
}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyID(id int) bool {
	err := db.QueryRow("SELECT * FROM movies WHERE id = ?", id)
	return err == nil 
}

