package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type Movie struct {
	ID 		int 		`json:"id"`
	Title	string		`json:"title"`
	TMDB_ID	int			`json:"tmdb_id"`
	Year	*int		`json:"release_year"`
	Genres	[]string	`json:"genre"`
	Status	string		`json:"status"`
}

// Filters for searching for movies
type Filters struct {
	Status 	*string	`json:"watched"`	// True means it has been watched
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
func getWatchlist(f *Filters) ([]Movie, error) {
	query := "SELECT * FROM movies"

	// Add requested filters
	selected_filters := []string{}
	if f.Status != nil {
		selected_filters = append(selected_filters, "status = \"" + *f.Status + "\"")
	}

	if len(selected_filters) != 0 {
		sql_filters := " WHERE " + strings.Join(selected_filters, " AND ")
		query += sql_filters
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Couldn't get movies from database")
		return nil, err
	}
	defer rows.Close()

	var watchlist []Movie

	for rows.Next() {
		var m Movie
		var genresJSON string

		err := rows.Scan(&m.ID, &m.TMDB_ID, &m.Title, &m.Year, &genresJSON, &m.Status)
		if err != nil {
			log.Println("Error reading movies from database!")
			return nil, err
		}

		// Process genres back into an array
		if genresJSON != "" {
			err = json.Unmarshal([]byte(genresJSON), &m.Genres)
		} else {
			m.Genres = []string{}
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
	genreBytes, err := json.Marshal(movie.Genres)
	if err != nil {
		log.Println("Couldn't convert the genres to JSON")
	}
	genreString := string(genreBytes)
	query := `
	INSERT INTO movies 
	(tmdb_id, title, release_year, genres)
	VALUES
	(?, ?, ?, ?)`
	//WHERE NOT EXISTS (
	//	SELECT 1 FROM movies WHERE tmdb_id = ?
	//)`
	_, err = db.Exec(query,
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

func removeMovie(id int) {
	query := "DELETE FROM movies WHERE id = ?"
	_, err := db.Exec(query, id)
	if (err != nil) {
		log.Printf("Failed to execute SQL", )
	}
}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyTMDB(tmdb_id int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM movies WHERE tmdb_id = ?)", tmdb_id).Scan(&exists)
	if err != nil {
		log.Printf("Couldn't query database for movie with tmdb_id: %d", tmdb_id)
		return false
	}
	return exists
}

// Check if a movie is in the database by tmdb ID
func isMovieInDBbyID(id int) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM movies WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		log.Printf("Couldn't query database for movie with id: %d", id)
		return false
	}
	return exists
}

