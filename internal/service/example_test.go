package service_test

import (
	"backend-test/internal/domain"
	"backend-test/internal/http/controller"
	"testing"
)

type MockBeerService struct {
	beers []domain.BeerStyle
}

func (m *MockBeerService) ListAllBeerStyles() ([]domain.BeerStyle, error) {
	return m.beers, nil
}

func (m *MockBeerService) GetBeerStyleByUUID(beerUUID string) (domain.BeerStyle, error) {
	for _, beer := range m.beers {
		if beer.UUID == beerUUID {
			return beer, nil
		}
	}
	return domain.BeerStyle{}, nil
}

func (m *MockBeerService) CreateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	m.beers = append(m.beers, beerStyle)
	return beerStyle, nil
}

func (m *MockBeerService) UpdateBeerStyle(beerStyle domain.BeerStyle) (domain.BeerStyle, error) {
	return beerStyle, nil
}

func (m *MockBeerService) DeleteBeerStyle(beerUUID string) error {
	return nil
}

type MockValidationService struct{}

func (m *MockValidationService) ValidateTemperatureRange(beerStyle domain.BeerStyle) error {
	return nil
}

func (m *MockValidationService) ValidateTemperatureInput(temperature float64) error {
	return nil
}

func (m *MockValidationService) ValidateUniqueNameForCreate(name string) error {
	return nil
}

func (m *MockValidationService) ValidateUniqueNameForUpdate(name string, excludeUUID string) error {
	return nil
}

func (m *MockValidationService) IsNoRowsError(err error) bool {
	return false
}

func (m *MockValidationService) ValidateUUID(uuidStr string) error {
	return nil
}

type MockUpdateService struct{}

func (m *MockUpdateService) ApplyBeerStyleUpdates(current *domain.BeerStyle, updates domain.BeerStyleUpdateRequest) bool {
	return false
}

func TestBeerControllerWithMocks(t *testing.T) {
	mockBeerService := &MockBeerService{
		beers: []domain.BeerStyle{
			{UUID: "123", Name: "Test Beer", TempMin: 0, TempMax: 5},
		},
	}
	mockValidationService := &MockValidationService{}
	mockUpdateService := &MockUpdateService{}

	beerController := controller.NewBeerController(
		mockBeerService,
		mockValidationService,
		mockUpdateService,
	)

	if beerController == nil {
		t.Error("BeerController should not be nil")
	}

}
