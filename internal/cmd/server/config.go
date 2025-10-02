package config

import (
	"os"

	"github.com/joho/godotenv"
)

var DATABASE_URL string

func init() {
	// Load environment variables from .env file if available
	_ = godotenv.Load()

	// Initialize configuration variables
	DATABASE_URL = GetDatabaseURL()
}

// GetDatabaseURL returns the database connection URL from environment or default
func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable"
	}
	return dbURL
}

// GetSpotifyClientID returns the Spotify API client ID from environment
func GetSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

// GetSpotifyClientSecret returns the Spotify API client secret from environment
func GetSpotifyClientSecret() string {
	return os.Getenv("SPOTIFY_CLIENT_SECRET")
}
