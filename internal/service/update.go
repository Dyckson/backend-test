package service

import (
	"backend-test/internal/domain"
)

type UpdateService struct{}

func NewUpdateService() *UpdateService {
	return &UpdateService{}
}

func (us *UpdateService) ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyle) bool {
	changed := false

	if updates.Name != "" && updates.Name != current.Name {
		current.Name = updates.Name
		changed = true
	}

	if updates.TempMin != 0 && updates.TempMin != current.TempMin {
		current.TempMin = updates.TempMin
		changed = true
	}

	if updates.TempMax != 0 && updates.TempMax != current.TempMax {
		current.TempMax = updates.TempMax
		changed = true
	}

	return changed
}

// GetChangedFields retorna os campos que foram alterados (para logs detalhados)
func (us *UpdateService) GetChangedFields(original, updated domain.BeerStyle) []string {
	var changedFields []string

	if updated.Name != "" && updated.Name != original.Name {
		changedFields = append(changedFields, "Name")
	}
	if updated.TempMin != 0 && updated.TempMin != original.TempMin {
		changedFields = append(changedFields, "TempMin")
	}
	if updated.TempMax != 0 && updated.TempMax != original.TempMax {
		changedFields = append(changedFields, "TempMax")
	}

	return changedFields
}
