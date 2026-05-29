package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/sevenbitscience/movies_to_watch/server/internal/config"
)

func getFileServer(w http.ResponseWriter, r *http.Request) {
	// Parse the webpage
	tmpl, err := template.ParseFiles(*config.Config.WebClientPath)
	if err != nil {
		http.Error(w, "Failed to parse web client", http.StatusInternalServerError)
	}

	dataURL := fmt.Sprintf("HTTP://%s", *config.Config.ServerURL)

	type PageData struct {
		DataURL string
	}

	data := PageData{
		DataURL: dataURL,
	}
	tmpl.Execute(w, data)
}
