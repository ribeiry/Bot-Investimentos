package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAssets_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
		{Ticker: "ITSA4", Market: "B3", Quantity: 200, AveragePrice: 10.00},
	}
	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)

	usecase := NewGetAssetUseCase(assetRepo)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BBSE3", result[0].Ticker)
	assetRepo.AssertExpectations(t)
}

func TestGetAssets_ListaVazia(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return([]domain.Asset{}, nil)

	usecase := NewGetAssetUseCase(assetRepo)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	assetRepo.AssertExpectations(t)
}

func TestGetAssets_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetAssetUseCase(assetRepo)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}
