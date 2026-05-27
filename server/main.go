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
	log.Println("Got TMDB API key")

	// Set up the DB
	initDB(os.Getenv("MOVIES_DATABASE_PATH"))
	log.Println("Connected to database")

	// Set up the REST API
	r := http.NewServeMux()

	r.HandleFunc("GET /movies", getMovies)
	r.HandleFunc("POST /movies/search", postSearchMovie)
	r.HandleFunc("POST /movies", postMovie)
	r.HandleFunc("PATCH /movies/{id}", patchMovie)

	server := http.Server{
		Addr:	os.Getenv("MOVIES_SERVER_URL"),
		Handler:	logging(r),
	}
	
	log.Printf("Server running on %s", server.Addr)
	server.ListenAndServe()
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path)
	})
}

// Format a reply to send some JSON
func sendJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get movies from the watchlist
// TODO Add paramaters to filter by status
func getMovies(w http.ResponseWriter, r *http.Request) {
	watchlist := getWatchlist()
	sendJSON(w, watchlist)
}

// searches for a movie and sends details to user, but don't save info.
func postSearchMovie(w http.ResponseWriter, r *http.Request) {
	var search struct {
		Query	string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		return	// Maybe let client know this failed??
	}

	movies := findMovies(search.Query)
	sendJSON(w, movies)
}


// Add a movie to the watchlist
func postMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie
	
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if (err != nil) {
		return	// Should probably let client know or something?
	}
	
	if (isMovieInDBbyTMDB(newMovie.TMDB_ID)) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	
	addMovie(&newMovie)
	w.WriteHeader(http.StatusCreated)
}



// Mark a movie as watched
func patchMovie(w http.ResponseWriter, r *http.Request) {
	var newStatus struct {
		Status	string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&newStatus)
	if (err != nil) {
		return
	}

	movieID, err := strconv.Atoi(r.PathValue("id"))
	if (err != nil || isMovieInDBbyID(movieID)) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	setWatchedStatus(movieID, newStatus.Status)
}

