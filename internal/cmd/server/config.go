package config

import (
	"os"

	"github.com/joho/godotenv"
)

var DATABASE_URL string

func init() {
	_ = godotenv.Load()

	DATABASE_URL = GetDatabaseURL()
}

func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable"
	}
	return dbURL
}

func GetSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func GetSpotifyClientSecret() string {
	return os.Getenv("SPOTIFY_CLIENT_SECRET")
}
