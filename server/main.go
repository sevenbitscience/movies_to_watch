package main

import (
	"encoding/json"
	"log"
	"net/http" // Handle RESTful API with just standard library
	"os"
	"strconv"

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
		log.Println(r.Method, r.URL.Path)
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
// TODO Add paramaters to filter by status
func getMovies(w http.ResponseWriter, r *http.Request) {
	// Get filters from the request, if there are any
	f := Filters{}
	if r.URL.Query().Has("status") {
		f.Status = new(string)
		*f.Status = r.URL.Query().Get("status")
	}

	watchlist, err := getWatchlist(&f)
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

	var movies []Movie
	findMovies(search.Query, &movies)

	sendJSON(w, movies)
}


// Add a movie to the watchlist
func postMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie
	
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Couldn't decode the JSON body")
		return
	}

	if (isMovieInDBbyTMDB(newMovie.TMDB_ID) == true) {
		log.Printf("Tried to add movie with duplicate TMDB ID: %d", newMovie.TMDB_ID)
		w.WriteHeader(http.StatusConflict)
		return
	}
	
	log.Printf("Adding a new movie %+v", newMovie)
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
		log.Println("Couldn't parse patch request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	movieID, err := strconv.Atoi(r.PathValue("id"))
	if (err != nil || !isMovieInDBbyID(movieID)) {
		log.Printf("Failed to find movie with ID %s", r.PathValue("id"))
		http.Error(w, "Couldn't find a movie with specified ID", http.StatusNotFound)
		return
	}

	log.Printf("Setting status of %d to %s", movieID, newStatus.Status)
	setWatchedStatus(movieID, newStatus.Status)
}

