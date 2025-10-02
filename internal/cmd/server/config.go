package config

import (
	"backend-test/external/spotify"
	"log"
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

func InitializeSpotifyService() *spotify.SpotifyService {
	clientID := GetSpotifyClientID()
	clientSecret := GetSpotifyClientSecret()

	if clientID == "" || clientSecret == "" {
		log.Println("Warning: Spotify credentials not set. Spotify integration will be disabled.")
		return nil
	}

	spotifyService, err := spotify.NewSpotifyService(clientID, clientSecret)
	if err != nil {
		log.Printf("Warning: Failed to initialize Spotify service: %v", err)
		return nil
	}

	return spotifyService
}
