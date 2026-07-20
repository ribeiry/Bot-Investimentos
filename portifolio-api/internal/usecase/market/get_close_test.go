package market

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClose_B3_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 40.00},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetCloseUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute("B3")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetClose_NYSE_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}
	quotes := []domain.Quote{
		{Ticker: "AAPL", CurrentValue: 160.00},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetCloseUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute("NYSE")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetClose_MarketInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	usecase := NewGetCloseUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute("INVALID")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestGetClose_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio").Return(nil, errors.New("db error"))

	usecase := NewGetCloseUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute("B3")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetClose_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetCloseUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute("B3")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}
