package market

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrices_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetPricesUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "BBSE3", result[0].Ticker)
	assert.Equal(t, 40.00, result[0].CurrentValue)

	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetPrices_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetPricesUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetPrices_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetPricesUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}
