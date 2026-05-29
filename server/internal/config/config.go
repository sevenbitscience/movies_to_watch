package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv" // Fetch info from .env file
)

var Config struct {
	ServerURL     *string
	TMBDkey       *string
	DatabasePath  *string
	WebAssetsPath *string
}

func InitConfig() error {
	if Config.ServerURL != nil {
		return errors.New("Config is already loaded")
	}

	// Try to load the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}

	// Load variables from environment or .env and return if any arent set
	serverURL, nf := os.LookupEnv("MOVIES_SERVER_URL"); if nf != true {
		return errors.New("'MOVIES_SERVER_URL' is not set") 
	}
	tmbdKey, nf := os.LookupEnv("TMDB_API_KEY"); if nf != true {
		return errors.New("'TMDB_API_KEY' is not set") 
	}
	databasePath, nf := os.LookupEnv("MOVIES_DATABASE_PATH"); if nf != true {
		return errors.New("'MOVIES_DATABASE_PATH' is not set") 
	}
	webAssetsPath, nf := os.LookupEnv("WEB_ASSETS_PATH"); if nf != true {
		return errors.New("'WEB_ASSETS_PATH' is not set") 
	}

	// Now set the data in the Config
	Config.ServerURL = new(string)
	*Config.ServerURL = serverURL

	Config.TMBDkey = new(string)
	*Config.TMBDkey = tmbdKey

	Config.DatabasePath = new(string)
	*Config.DatabasePath = databasePath

	Config.WebAssetsPath = new(string)
	*Config.WebAssetsPath = webAssetsPath

	return nil
}
