package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSummary_Realtime_Success(t *testing.T) {
	// 1. Cria os mocks
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)

	// 2. Define o comportamento esperado
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 40.00},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	// 3. Executa o usecase
	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute("realtime")

	// 4. Verifica o resultado
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

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 39.00},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	cachedProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute("cached")

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

	assetRepo.On("ReturnAllPortfolio").Return(nil, errors.New("db error"))

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute("realtime")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetSummary_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	cachedProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
	}

	assetRepo.On("ReturnAllPortfolio").Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetSummaryUseCase(assetRepo, marketProvider, cachedProvider)
	result, err := usecase.Execute("realtime")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}
