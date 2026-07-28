package portifolio

import (
	"database/sql"
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var _ = time.Now() // força import de time

func TestGetPeriodSummary_Weekly_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 42.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "weekly", result.Period)
	assert.Len(t, result.Assets, 1)
	assert.Equal(t, "BBSE3", result.Assets[0].Ticker)
	assert.Equal(t, 40.00, result.Assets[0].PriceStart)
	assert.Equal(t, 42.00, result.Assets[0].PriceCurrent)
	assert.InDelta(t, 200.00, result.Assets[0].ChangeValue, 0.01)
	assert.InDelta(t, 5.0, result.Assets[0].ChangePercent, 0.01)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}

func TestGetPeriodSummary_Monthly_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}
	quotes := []domain.Quote{
		{Ticker: "AAPL", CurrentValue: 160.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	priceHistory.On("GetPriceAtDate", "AAPL", mock.AnythingOfType("time.Time")).Return(155.00, nil)

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "monthly")

	assert.NoError(t, err)
	assert.Equal(t, "monthly", result.Period)
	assert.InDelta(t, 50.00, result.TotalChangeValue, 0.01)
	assetRepo.AssertExpectations(t)
}

func TestGetPeriodSummary_SemHistorico_UsaAveragePrice(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 42.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(0.0, sql.ErrNoRows)

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	assert.Equal(t, 38.00, result.Assets[0].PriceStart) // usa average_price como fallback
	assert.Equal(t, 42.00, result.Assets[0].PriceCurrent)
	assetRepo.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}

func TestGetPeriodSummary_PeriodoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "yearly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestGetPeriodSummary_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetPeriodSummary_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetPeriodSummary_MultiploAtivos(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 40.00},
		{Ticker: "AAPL", CurrentValue: 160.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(38.00, nil)
	priceHistory.On("GetPriceAtDate", "AAPL", mock.AnythingOfType("time.Time")).Return(150.00, nil)

	usecase := NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	assert.Len(t, result.Assets, 2)
	// BBSE3: 100*(40-38) = 200, AAPL: 10*(160-150) = 100 → total = 300
	assert.InDelta(t, 300.00, result.TotalChangeValue, 0.01)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}
