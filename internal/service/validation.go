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

// ValidateTemperatureRange valida se TempMin < TempMax e se estão em faixas razoáveis
func (vs *ValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	// Valida limites baseados nas temperaturas extremas da Terra
	// Menor: -89,2°C (Antártida) | Maior: +56,7°C (Vale da Morte)
	// Usando range ampliado para segurança: -90°C a +60°C
	if beerStyle.TempMin < -90 || beerStyle.TempMin > 60 {
		return fmt.Errorf("minimum temperature (%.1f) must be between -90°C and 60°C", beerStyle.TempMin)
	}

	if beerStyle.TempMax < -90 || beerStyle.TempMax > 60 {
		return fmt.Errorf("maximum temperature (%.1f) must be between -90°C and 60°C", beerStyle.TempMax)
	}

	// Valida se TempMin < TempMax
	if beerStyle.TempMin >= beerStyle.TempMax {
		return fmt.Errorf("minimum temperature (%.1f) must be less than maximum temperature (%.1f)",
			beerStyle.TempMin, beerStyle.TempMax)
	}

	return nil
}

// ValidateTemperatureInput valida se a temperatura de entrada está dentro dos limites aceitáveis
// Esta validação é usada para entrada de temperatura nas APIs de recomendação
func (vs *ValidationService) ValidateTemperatureInput(temperature float64) error {
	// Valida limites baseados nas temperaturas extremas da Terra
	// Menor: -89,2°C (Antártida) | Maior: +56,7°C (Vale da Morte)
	// Usando range ampliado para segurança: -90°C a +60°C
	if temperature < -90 || temperature > 60 {
		return fmt.Errorf("temperature (%.1f) must be between -90°C and 60°C", temperature)
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
