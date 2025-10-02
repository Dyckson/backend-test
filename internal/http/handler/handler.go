package handler

import (
	"backend-test/external/spotify"
	"backend-test/internal/http/controller"
	"backend-test/internal/service"
	"backend-test/internal/storage/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var beerController *controller.BeerController
var recommendationController *controller.RecommendationController

func init() {
	// Inicializa as dependências
	beerRepo := repository.BeerRepository{}
	beerService := service.NewBeerService(beerRepo)
	validationService := service.NewValidationService(*beerService)
	updateService := service.NewUpdateService()

	// Inicializa Spotify Service
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Println("Warning: Spotify credentials not set. Using placeholder.")
		clientID = "your_client_id"
		clientSecret = "your_client_secret"
	}

	spotifyService, err := spotify.NewSpotifyService(clientID, clientSecret)
	if err != nil {
		log.Printf("Warning: Failed to initialize Spotify service: %v", err)
		spotifyService = nil
	}

	recommendationService := service.NewRecommendationService(*beerService, spotifyService)

	// Inicializa controllers
	beerController = controller.NewBeerController(*beerService, *validationService, *updateService)
	recommendationController = controller.NewRecommendationController(recommendationService, validationService)
}

func HealthCheckStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK.",
	})
}

func HandleRequests(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/check", HealthCheckStatus)

	beer := api.Group("/beer-styles")
	beer.GET("/list", beerController.ListAllBeerStyles)
	beer.POST("/create", beerController.CreateBeerStyle)
	beer.PUT("/edit/:beerUUID", beerController.UpdateBeerStyle)
	beer.DELETE("/:beerUUID", beerController.DeleteBeerStyle)

	recommendations := api.Group("/recommendations")
	recommendations.POST("/suggest", recommendationController.SuggestSpotifyPlaylist)
}
