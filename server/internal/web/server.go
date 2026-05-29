package web

import "net/http"

func SetupServer(url string) *http.Server {
	r := setupRoutes()

	server := http.Server{
		Addr:    url,
		Handler: corsMiddleware(logging(r)),
	}

	return &server
}
