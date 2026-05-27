package main

import (
	"os"
	"log"
	"strconv"
	"net/http"

	"github.com/gin-gonic/gin" // Handle the RESTful API
	"github.com/joho/godotenv" // Fetch info from .env file
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

	// Set up Gin to handle the REST API
	r := gin.Default()

	r.GET("/movies", getMovies)
	r.POST("/movies/search", postSearchMovie)
	r.POST("/movies", postMovie)
	r.PATCH("/movies/:id", patchMovie)

	r.Run(os.Getenv("MOVIES_SERVER_URL"))
}

// get movies from the watchlist
// TODO Add paramaters to filter by status
func getMovies(c *gin.Context) {
	watchlist := getWatchlist()
	c.IndentedJSON(http.StatusOK, watchlist)
}

type SearchMovie struct {
	Query	string `json:"query"`
}

// searches for a movie and sends details to user, but don't save info.
func postSearchMovie(c *gin.Context) {
	var search SearchMovie

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


type MovieStatus struct {
	Status	string `json:"status"`
}

// Mark a movie as watched
func patchMovie(c *gin.Context) {
	var newStatus MovieStatus

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
