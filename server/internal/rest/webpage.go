package rest

import (
	"net/http"

	"github.com/sevenbitscience/movies_to_watch/server/internal/config"
)

func getFileServer() http.Handler {
	files := http.FileServer(http.Dir(*config.Config.WebAssetsPath))
	return files
}
	

