package handler

import (
	"backend-test/internal/http/controller"
	"backend-test/internal/service"
	"backend-test/internal/storage/repository"

	"github.com/gin-gonic/gin"
)

var beerController *controller.BeerController

func init() {
	// Inicializa as dependências
	beerRepo := repository.BeerRepository{}
	beerService := service.NewBeerService(beerRepo)
	validationService := service.NewValidationService(*beerService)
	updateService := service.NewUpdateService()
	beerController = controller.NewBeerController(*beerService, *validationService, *updateService)
}

func HandleRequests(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/check", controller.HealthCheckStatus)

	beer := api.Group("/beer-styles")
	beer.GET("/list", beerController.ListAllBeerStyles)

	beer.POST("/create", beerController.CreateBeerStyle)
	//beer.POST("/suggest", beerController.SuggestSpotifyPlaylist)

	beer.PUT("/edit/:beerUUID", beerController.UpdateBeerStyle)

	beer.DELETE("/:beerUUID", beerController.DeleteBeerStyle)
}
