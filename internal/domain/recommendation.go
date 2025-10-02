package domain

// DTOs for the Spotify recommendation feature

// TemperatureRequest represents the request body for temperature-based recommendations
type TemperatureRequest struct {
	Temperature float64 `json:"temperature"`
}

// TrackInfo represents a single track from a Spotify playlist
type TrackInfo struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
	Link   string `json:"link"`
}

// PlaylistInfo represents a Spotify playlist with its tracks
type PlaylistInfo struct {
	Name   string      `json:"name"`
	Tracks []TrackInfo `json:"tracks"`
}

// RecommendationResponse represents the complete recommendation response
// containing the best beer style and corresponding Spotify playlist
type RecommendationResponse struct {
	BeerStyle string       `json:"beerStyle"`
	Playlist  PlaylistInfo `json:"playlist"`
}
