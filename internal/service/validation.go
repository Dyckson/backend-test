package service

import (
	"backend-test/internal/domain"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type ValidationService struct {
	beerService BeerService
}

func NewValidationService(beerService BeerService) *ValidationService {
	return &ValidationService{
		beerService: beerService,
	}
}

// isNoRowsError verifica se o erro indica que não há linhas no resultado
func (vs *ValidationService) isNoRowsError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set")
}

// IsNoRowsError método público para verificar se é erro "no rows"
func (vs *ValidationService) IsNoRowsError(err error) bool {
	return vs.isNoRowsError(err)
}

// ValidateTemperatureRange valida se TempMin < TempMax
func (vs *ValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	if beerStyle.TempMin >= beerStyle.TempMax {
		return fmt.Errorf("minimum temperature (%.1f) must be less than maximum temperature (%.1f)",
			beerStyle.TempMin, beerStyle.TempMax)
	}
	return nil
}

// ValidateUniqueNameForCreate valida se o nome da cerveja é único para criação
func (vs *ValidationService) ValidateUniqueNameForCreate(name string) error {
	beerStyles, err := vs.beerService.ListAllBeerStyles()
	if err != nil {
		if !vs.isNoRowsError(err) {
			return fmt.Errorf("failed to check beer styles: %w", err)
		}
		// Se não há estilos, o nome é único
		return nil
	}

	for _, style := range beerStyles {
		if style.Name == name {
			return fmt.Errorf("beer style with name '%s' already exists", name)
		}
	}

	return nil
}

// ValidateUniqueNameForUpdate valida se o nome da cerveja é único para atualização
func (vs *ValidationService) ValidateUniqueNameForUpdate(name string, excludeUUID string) error {
	if name == "" {
		return nil // Nome vazio não precisa validar
	}

	beerStyles, err := vs.beerService.ListAllBeerStyles()
	if err != nil {
		if !vs.isNoRowsError(err) {
			return fmt.Errorf("failed to check beer styles: %w", err)
		}
		// Se não há estilos, o nome é único
		return nil
	}

	for _, style := range beerStyles {
		if style.Name == name && style.UUID != excludeUUID {
			return fmt.Errorf("beer style with name '%s' already exists", name)
		}
	}

	return nil
}
