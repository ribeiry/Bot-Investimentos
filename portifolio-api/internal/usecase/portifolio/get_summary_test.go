package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSummary_Realtime_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute(userID, "realtime")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.MarketSummary, 1)
	assert.Equal(t, "B3", result.MarketSummary[0].MarketProvider)

	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetSummary_Cached_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 39.00}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	cachedProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute(userID, "cached")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.MarketSummary, 1)
	assert.Equal(t, "B3", result.MarketSummary[0].MarketProvider)

	assetRepo.AssertExpectations(t)
	cachedProvider.AssertExpectations(t)
	marketProvider.AssertNotCalled(t, "GetByTickers")
}

func TestGetSummary_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute(userID, "realtime")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetSummary_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute(userID, "realtime")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetSummary_CachedProvider_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	cachedProvider.On("GetByTickers", assets).Return(nil, errors.New("cache error"))

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute(userID, "cached")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	cachedProvider.AssertExpectations(t)
}
