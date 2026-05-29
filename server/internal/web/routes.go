package web

import (
	"net/http"
)

func setupRoutes() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/", getFileServer)
	r.HandleFunc("GET /movies", getMovies)
	r.HandleFunc("POST /movies/search", postSearchMovie)
	r.HandleFunc("POST /movies", postMovie)
	r.HandleFunc("PATCH /movies/{id}", patchMovie)
	r.HandleFunc("DELETE /movies/{id}", deleteMovie)

	return r
}
