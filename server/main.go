package main

import (
	"encoding/json"
	"log"
	"net/http" // Handle RESTful API with just standard library
	"os"
	"strconv"
	"strings"

	"github.com/sevenbitscience/movies_to_watch/server/internal/types"
	"github.com/sevenbitscience/movies_to_watch/server/internal/database"
	"github.com/sevenbitscience/movies_to_watch/server/internal/tmdb"

	"github.com/joho/godotenv" // Fetch info from .env file
)

func main() {
	// Fetch arguments from the .env file
	if err := godotenv.Load()
	err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Pass the TMDB API key over to tmdb.go
	tmdb.SetApiKey(os.Getenv("TMDB_API_KEY"))
	log.Println("Got TMDB API key")

	// Set up the DB
	database.InitDB(os.Getenv("MOVIES_DATABASE_PATH"))
	log.Println("Connected to database")

	// Set up the REST API
	r := http.NewServeMux()

	r.HandleFunc("GET /movies", getMovies)
	r.HandleFunc("POST /movies/search", postSearchMovie)
	r.HandleFunc("POST /movies", postMovie)
	r.HandleFunc("PATCH /movies/{id}", patchMovie)
	r.HandleFunc("DELETE /movies/{id}", deleteMovie)

	server := http.Server{
		Addr:	os.Getenv("MOVIES_SERVER_URL"),
		Handler:	corsMiddleware(logging(r)),
	}
	
	log.Printf("Server running on %s", server.Addr)
	server.ListenAndServe()
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        // If it's a preflight OPTIONS request, stop here and reply with 204
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
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
func getMovies(w http.ResponseWriter, r *http.Request) {
	// Get filters from the request, if there are any
	f := database.Filters{}
	// Filter by status
	if r.URL.Query().Has("status") {
		f.Status = new(string)
		*f.Status = r.URL.Query().Get("status")
	}
	// Filter by Genre
	if r.URL.Query().Has("genres") {
		f.Genres = strings.Split(r.URL.Query().Get("genres"), ",")
	}

	watchlist, err := database.GetWatchlist(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError);
		return
	}
	if len(watchlist) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	sendJSON(w, watchlist)
}

// searches for a movie and sends details to user, but don't save info.
func postSearchMovie(w http.ResponseWriter, r *http.Request) {
	var search struct {
		Query	string `json:"query"`
	}

	err := json.NewDecoder(r.Body).Decode(&search)
	if err != nil {
		log.Printf("Couldn't parse search term")
		return	// Maybe let client know this failed??
	}

	var movies []types.Movie
	tmdb.FindMovies(search.Query, &movies)

	sendJSON(w, movies)
}


// Add a movie to the watchlist
func postMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie types.Movie
	
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't decode the JSON body")
		return
	}

	if (database.IsMovieInDBbyTMDB(newMovie.TMDB_ID) == true) {
		log.Printf("Tried to add movie with duplicate TMDB ID: %d", newMovie.TMDB_ID)
		w.WriteHeader(http.StatusConflict)
		return
	}
	
	log.Printf("Adding a new movie %+v", newMovie)
	database.AddMovie(&newMovie)
	w.WriteHeader(http.StatusCreated)
}



// Mark a movie as watched
func patchMovie(w http.ResponseWriter, r *http.Request) {
	var newStatus struct {
		Status	string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&newStatus)
	if (err != nil) {
		log.Println("Couldn't parse patch request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	movieID, err := strconv.Atoi(r.PathValue("id"))
	if (err != nil || !database.IsMovieInDBbyID(movieID)) {
		log.Printf("Failed to find movie with ID %s", r.PathValue("id"))
		http.Error(w, "Couldn't find a movie with specified ID", http.StatusNotFound)
		return
	}

	log.Printf("Setting status of %d to %s", movieID, newStatus.Status)
	database.SetWatchedStatus(movieID, newStatus.Status)
}


func deleteMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("id"))
	if (err != nil || !database.IsMovieInDBbyID(movieID)) {
		log.Printf("Failed to find movie with ID %s", r.PathValue("id"))
		http.Error(w, "Couldn't find a movie with specified ID", http.StatusNotFound)
		return
	}
	log.Printf("Deleting movie %d", movieID)
	database.RemoveMovie(movieID)
}

