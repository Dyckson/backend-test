package controller

import (
	"backend-test/internal/domain"
	"backend-test/internal/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BeerController struct {
	BeerService       service.BeerService
	ValidationService service.ValidationService
	UpdateService     service.UpdateService
}

func NewBeerController(beerService service.BeerService, validationService service.ValidationService, updateService service.UpdateService) *BeerController {
	return &BeerController{
		BeerService:       beerService,
		ValidationService: validationService,
		UpdateService:     updateService,
	}
}

func (bc *BeerController) ListAllBeerStyles(c *gin.Context) {
	beerStyles, err := bc.BeerService.ListAllBeerStyles()
	if err != nil {
		log.Printf("controller=BeerController func=ListAllBeerStyles err=%v", err)

		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "no beer styles found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"beerStyles": beerStyles,
	})
}

func (bc *BeerController) CreateBeerStyle(c *gin.Context) {
	var inputStyle domain.BeerStyle
	if err := c.ShouldBindJSON(&inputStyle); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body, name, temp_min and temp_max are required",
		})
		return
	}

	// Valida se o nome é único
	if err := bc.ValidationService.ValidateUniqueNameForCreate(inputStyle.Name); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)

		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to validate beer style name",
		})
		return
	}

	// Valida a faixa de temperatura
	if err := bc.ValidationService.ValidateTemperatureRange(inputStyle); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newBeerStyle, err := bc.BeerService.CreateBeerStyle(inputStyle)
	if err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create beer style",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newBeerStyle,
	})
}

func (bc *BeerController) UpdateBeerStyle(c *gin.Context) {
	beerUUID := c.Param("beerUUID")
	if beerUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "beerUUID is required",
		})
		return
	}

	var updateRequest domain.BeerStyleUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	previewBeerStyle := domain.BeerStyle{
		Name:    updateRequest.Name,
		TempMin: updateRequest.TempMin,
		TempMax: updateRequest.TempMax,
	}

	currentBeerStyle, err := bc.BeerService.GetBeerStyleByUUID(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "beer style not found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	if previewBeerStyle.Name != "" && previewBeerStyle.Name != currentBeerStyle.Name {
		if err := bc.ValidationService.ValidateUniqueNameForUpdate(previewBeerStyle.Name, currentBeerStyle.UUID); err != nil {
			log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)

			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to validate beer style name",
			})
			return
		}
	}

	// Aplica as mudanças usando o UpdateService
	originalStyle := currentBeerStyle // Cópia para comparação
	changed := bc.UpdateService.ApplyBeerStyleUpdates(&currentBeerStyle, previewBeerStyle)

	if changed {
		// Valida a faixa de temperatura após as mudanças
		if err := bc.ValidationService.ValidateTemperatureRange(currentBeerStyle); err != nil {
			log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}

		// Log dos campos alterados para auditoria
		changedFields := bc.UpdateService.GetChangedFields(originalStyle, previewBeerStyle)
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s changed_fields=%v", beerUUID, changedFields)
	}

	if !changed {
		c.JSON(http.StatusOK, gin.H{
			"message": "No changes detected.",
			"data":    currentBeerStyle,
		})
		return
	}

	updatedBeerStyle, err := bc.BeerService.UpdateBeerStyle(currentBeerStyle)
	if err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update beer style",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beer style updated.",
		"data":    updatedBeerStyle,
	})
}

func (bc *BeerController) DeleteBeerStyle(c *gin.Context) {
	beerUUID := c.Param("beerUUID")
	if beerUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "beerUUID is required",
		})
		return
	}

	_, err := bc.BeerService.GetBeerStyleByUUID(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=DeleteBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "beer style not found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	err = bc.BeerService.DeleteBeerStyle(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=DeleteBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beer style deleted successfully",
	})
}
