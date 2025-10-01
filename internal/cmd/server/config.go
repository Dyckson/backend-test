package config

import "os"

func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable"
	}
	return dbURL
}

var DATABASE_URL = GetDatabaseURL()
