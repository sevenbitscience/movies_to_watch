package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http" // Handle RESTful API with just standard library
	"strconv"
	"strings"

	"github.com/sevenbitscience/movies_to_watch/server/internal/database"
	"github.com/sevenbitscience/movies_to_watch/server/internal/tmdb"
	"github.com/sevenbitscience/movies_to_watch/server/internal/types"
)

// Format a reply to send some JSON
func encodeBody[T any](w http.ResponseWriter, data T) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func decodeBody[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return v, fmt.Errorf("Failed decoding JSON %w", err)
	}
	return v, nil
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(watchlist) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	encodeBody(w, watchlist)
}

// searches for a movie and sends details to user, but don't save info.
func postSearchMovie(w http.ResponseWriter, r *http.Request) {
	type SearchRequest struct {
		Query string `json:"query"`
	}

	search, err := decodeBody[SearchRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var movies []types.Movie
	tmdb.FindMovies(search.Query, &movies)

	encodeBody(w, movies)
}

// Add a movie to the watchlist
func postMovie(w http.ResponseWriter, r *http.Request) {
	newMovie, err := decodeBody[types.Movie](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if database.IsMovieInDBbyTMDB(newMovie.TMDB_ID) == true {
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
	type StatusChangeRequest struct {
		Status string `json:"status"`
	}

	newStatus, err := decodeBody[StatusChangeRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || !database.IsMovieInDBbyID(movieID) {
		log.Printf("Failed to find movie with ID %s", r.PathValue("id"))
		http.Error(w, "Couldn't find a movie with specified ID", http.StatusNotFound)
		return
	}

	log.Printf("Setting status of %d to %s", movieID, newStatus.Status)
	database.SetWatchedStatus(movieID, newStatus.Status)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	movieID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || !database.IsMovieInDBbyID(movieID) {
		log.Printf("Failed to find movie with ID %s", r.PathValue("id"))
		http.Error(w, "Couldn't find a movie with specified ID", http.StatusNotFound)
		return
	}

	log.Printf("Deleting movie %d", movieID)
	database.RemoveMovie(movieID)
}
