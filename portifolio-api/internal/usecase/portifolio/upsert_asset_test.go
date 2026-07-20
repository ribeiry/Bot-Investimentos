package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsertAsset_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "BBSE3",
		Market:       "B3",
		Quantity:     100,
		AveragePrice: 38.50,
	}

	assetRepo.On("Upsert", asset).Return(nil)

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.NoError(t, err)
	assetRepo.AssertExpectations(t)
}

func TestUpsertAsset_TickerVazio(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "",
		Market:       "B3",
		Quantity:     100,
		AveragePrice: 38.50,
	}

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAsset_QuantidadeInvalida(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "BBSE3",
		Market:       "B3",
		Quantity:     0,
		AveragePrice: 38.50,
	}

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAsset_PrecoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "BBSE3",
		Market:       "B3",
		Quantity:     100,
		AveragePrice: 0,
	}

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAsset_MercadoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "BBSE3",
		Market:       "INVALID",
		Quantity:     100,
		AveragePrice: 38.50,
	}

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.Error(t, err)
	assetRepo.AssertExpectations(t)
}
func TestUpsertAsset_RepError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	asset := domain.Asset{
		Ticker:       "BBSE3",
		Market:       "B3",
		Quantity:     100,
		AveragePrice: 38.50,
	}

	assetRepo.On("Upsert", asset).Return(errors.New("db error"))

	usecase := NewUpsertAssetUseCase(assetRepo)
	err := usecase.Execute(asset)

	assert.Error(t, err)
	assetRepo.AssertExpectations(t)
} // ← garante que esse } está aqui
