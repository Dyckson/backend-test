package service

import (
	"backend-test/external/spotify"
	"backend-test/internal/domain"
	"fmt"
	"log"

	spotifyapi "github.com/zmb3/spotify/v2"
)

type RecommendationService struct {
	beerService    BeerService
	spotifyService *spotify.SpotifyService
}

func NewRecommendationService(beerService BeerService, spotifyService *spotify.SpotifyService) *RecommendationService {
	return &RecommendationService{
		beerService:    beerService,
		spotifyService: spotifyService,
	}
}

// FindBestBeerStyleForTemperature encontra o estilo de cerveja mais adequado para uma temperatura
// Regra: seleciona o estilo cuja média de temperaturas (TempMin + TempMax)/2 está mais próxima do input
// Em caso de empate, ordena alfabeticamente
func (rs *RecommendationService) FindBestBeerStyleForTemperature(temperature float64) (*domain.BeerStyle, error) {
	// Busca todos os estilos de cerveja
	allBeerStyles, err := rs.beerService.ListAllBeerStyles()
	if err != nil {
		return nil, fmt.Errorf("failed to get beer styles: %w", err)
	}

	if len(allBeerStyles) == 0 {
		return nil, fmt.Errorf("no beer styles found")
	}

	type candidateStyle struct {
		style    *domain.BeerStyle
		distance float64
		average  float64
	}

	var candidates []candidateStyle
	var minDistance float64 = 999999

	// Calcula a distância da temperatura input para a média de cada estilo
	for i := range allBeerStyles {
		beerStyle := &allBeerStyles[i]
		average := (beerStyle.TempMin + beerStyle.TempMax) / 2
		distance := abs(temperature - average)
		// log para debug
		// log.Printf("BeerStyle: %s, Range: [%.1f, %.1f], Average: %.1f, Distance: %.1f", beerStyle.Name, beerStyle.TempMin, beerStyle.TempMax, average, distance)

		// Se encontrou uma distância menor, reinicia a lista de candidatos
		if distance < minDistance {
			minDistance = distance
			candidates = []candidateStyle{{
				style:    beerStyle,
				distance: distance,
				average:  average,
			}}
		} else if distance == minDistance {
			// Se a distância é igual, adiciona à lista de candidatos
			candidates = append(candidates, candidateStyle{
				style:    beerStyle,
				distance: distance,
				average:  average,
			})
		}
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no suitable beer style found for temperature %.1f°C", temperature)
	}

	// Se há apenas um candidato, retorna ele
	if len(candidates) == 1 {
		return candidates[0].style, nil
	}

	// Se há múltiplos candidatos com a mesma distância, ordena alfabeticamente
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].style.Name > candidates[j].style.Name {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	bestMatch := candidates[0].style
	return bestMatch, nil
} // GetRecommendationForTemperature retorna recomendação completa com cerveja e playlist
func (rs *RecommendationService) GetRecommendationForTemperature(temperature float64) (*domain.RecommendationResponse, error) {
	// 1. Encontra o melhor estilo de cerveja
	beerStyle, err := rs.FindBestBeerStyleForTemperature(temperature)
	if err != nil {
		return nil, err
	}

	// 2. Busca playlist no Spotify (se disponível)
	var tracks []domain.TrackInfo
	var playlistName string

	if rs.spotifyService != nil {
		playlist, err := rs.spotifyService.SearchPlaylistByName(beerStyle.Name)
		if err != nil {
			log.Printf("Failed to find Spotify playlist for %s: %v", beerStyle.Name, err)
			// Retorna erro específico quando não encontra playlist
			return nil, fmt.Errorf("no playlist found for beer style '%s'", beerStyle.Name)
		}

		playlistName = playlist.Name
		tracks = rs.convertSpotifyTracks(playlist)

		// Se a playlist não tem tracks, também consideramos como não encontrada
		if len(tracks) == 0 {
			return nil, fmt.Errorf("playlist '%s' found but contains no valid tracks", playlistName)
		}
	} else {
		// Spotify não disponível - retorna erro em produção
		// Em desenvolvimento, pode retornar dados de exemplo
		log.Println("Spotify service not available")
		return nil, fmt.Errorf("spotify service unavailable")
	}

	response := &domain.RecommendationResponse{
		BeerStyle: beerStyle.Name,
		Playlist: domain.PlaylistInfo{
			Name:   playlistName,
			Tracks: tracks,
		},
	}

	return response, nil
}

// convertSpotifyTracks converte tracks do Spotify para nosso formato
func (rs *RecommendationService) convertSpotifyTracks(playlist *spotifyapi.FullPlaylist) []domain.TrackInfo {
	tracks := make([]domain.TrackInfo, 0)
	if len(playlist.Tracks.Tracks) > 0 {
		// Limita a 10 tracks para não sobrecarregar a resposta
		maxTracks := len(playlist.Tracks.Tracks)
		if maxTracks > 10 {
			maxTracks = 10
		}

		for i := 0; i < maxTracks; i++ {
			track := playlist.Tracks.Tracks[i].Track
			if track.Name != "" {
				// Monta o link do Spotify
				spotifyLink := fmt.Sprintf("https://open.spotify.com/track/%s", track.ID)

				// Pega o primeiro artista (se houver)
				artistName := "Unknown Artist"
				if len(track.Artists) > 0 {
					artistName = track.Artists[0].Name
				}

				tracks = append(tracks, domain.TrackInfo{
					Name:   track.Name,
					Artist: artistName,
					Link:   spotifyLink,
				})
			}
		}
	}
	return tracks
}

// abs retorna o valor absoluto de um float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
