package main

import (
	"os"
	"log"
	"encoding/json"
	"strconv"
	"net/http"	// Handle RESTful API with just standard library

	"github.com/joho/godotenv"	// Fetch info from .env file
)

func main() {
	// Fetch arguments from the .env file
	if err := godotenv.Load()
	err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Pass the TMDB API key over to tmdb.go
	setApiKey(os.Getenv("TMDB_API_KEY"))

	// Set up the DB
	initDB(os.Getenv("MOVIES_DATABASE_PATH"))

	// Set up the REST API
	r := http.NewServeMux()

	r.HandleFunc("GET /movies", getMovies)
	r.HandleFunc("POST /movies/search", postSearchMovie)
	r.HandleFunc("POST /movies", postMovie)
	r.HandleFunc("PATCH /movies/:id", patchMovie)

	server := http.Server{
		Addr:	os.Getenv("MOVIES_SERVER_URL"),
		Handler:	r,
	}
	server.ListenAndServe()
}

// Format a reply to send some JSON
func sendJSON(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	return err
}

// get movies from the watchlist
// TODO Add paramaters to filter by status
func getMovies(w http.ResponseWriter, r *http.Request) {
	watchlist := getWatchlist()
	sendJSON(w, watchlist)
}

// searches for a movie and sends details to user, but don't save info.
func postSearchMovie(c *gin.Context) {
	var search struct {
		Query	string `json:"query"`
	}

	err := c.ShouldBindJSON(&search)
	if err != nil {
		return	// Maybe let client know this failed??
	}

	movies := findMovies(search.Query)
	c.IndentedJSON(http.StatusOK, movies)
}


// Add a movie to the watchlist
func postMovie(c *gin.Context) {
	var newMovie Movie
	
	err := c.ShouldBindJSON(&newMovie)
	if (err != nil) {
		return	// Should probably let client know or something?
	}
	
	if (isMovieInDBbyTMDB(newMovie.TMDB_ID)) {
		c.Status(http.StatusConflict)
		return
	}
	
	addMovie(&newMovie)
	c.Status(http.StatusCreated)
}



// Mark a movie as watched
func patchMovie(c *gin.Context) {
	var newStatus struct {
		Status	string `json:"status"`
	}

	err := c.ShouldBindJSON(&newStatus)
	if (err != nil) {
		return
	}

	movieID, err := strconv.Atoi(c.Param("id"))
	if (err != nil || isMovieInDBbyID(movieID)) {
		c.Status(http.StatusNotFound)
		return
	}

	setWatchedStatus(movieID, newStatus.Status)
}
