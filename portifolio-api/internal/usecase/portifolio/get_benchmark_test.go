package portifolio

import (
	"database/sql"
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var benchAssets = []domain.Asset{
	{Ticker: "^BVSP", Market: "B3"},
	{Ticker: "SPX", Market: "NYSE"},
}

func TestGetBenchmark_Weekly_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	portfolioQuotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 42.00},
	}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000.00},
		{Ticker: "SPX", CurrentValue: 5500.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)

	// portfolio history
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)
	// benchmark history
	priceHistory.On("GetPriceAtDate", "^BVSP", mock.AnythingOfType("time.Time")).Return(125000.00, nil)
	priceHistory.On("GetPriceAtDate", "SPX", mock.AnythingOfType("time.Time")).Return(5000.00, nil)

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "weekly", result.Period)
	// portfolio: (4200-4000)/4000 * 100 = 5%
	assert.InDelta(t, 5.0, result.PortfolioReturn, 0.01)
	assert.Len(t, result.Benchmarks, 2)

	ibov := result.Benchmarks[0]
	assert.Equal(t, "IBOV", ibov.Name)
	// ibov: (130000-125000)/125000 * 100 = 4%
	assert.InDelta(t, 4.0, ibov.ReturnPercent, 0.01)
	// relative: 5 - 4 = 1
	assert.InDelta(t, 1.0, ibov.RelativePerformance, 0.01)

	sp500 := result.Benchmarks[1]
	assert.Equal(t, "S&P500", sp500.Name)
	// sp500: (5500-5000)/5000 * 100 = 10%
	assert.InDelta(t, 10.0, sp500.ReturnPercent, 0.01)
	// relative: 5 - 10 = -5
	assert.InDelta(t, -5.0, sp500.RelativePerformance, 0.01)

	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}

func TestGetBenchmark_Monthly_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}
	portfolioQuotes := []domain.Quote{{Ticker: "AAPL", CurrentValue: 165.00}}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000.00},
		{Ticker: "SPX", CurrentValue: 5500.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)
	priceHistory.On("GetPriceAtDate", "AAPL", mock.AnythingOfType("time.Time")).Return(150.00, nil)
	priceHistory.On("GetPriceAtDate", "^BVSP", mock.AnythingOfType("time.Time")).Return(125000.00, nil)
	priceHistory.On("GetPriceAtDate", "SPX", mock.AnythingOfType("time.Time")).Return(5000.00, nil)

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "monthly")

	assert.NoError(t, err)
	assert.Equal(t, "monthly", result.Period)
	// AAPL: (1650-1500)/1500*100 = 10%
	assert.InDelta(t, 10.0, result.PortfolioReturn, 0.01)
	assetRepo.AssertExpectations(t)
}

func TestGetBenchmark_SemHistoricoBenchmark_NoHistory(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	portfolioQuotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 42.00}}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000.00},
		{Ticker: "SPX", CurrentValue: 5500.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)
	priceHistory.On("GetPriceAtDate", "^BVSP", mock.AnythingOfType("time.Time")).Return(0.0, sql.ErrNoRows)
	priceHistory.On("GetPriceAtDate", "SPX", mock.AnythingOfType("time.Time")).Return(0.0, sql.ErrNoRows)

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	assert.Len(t, result.Benchmarks, 2)
	assert.True(t, result.Benchmarks[0].NoHistory)
	assert.True(t, result.Benchmarks[1].NoHistory)
	assert.Equal(t, 130000.00, result.Benchmarks[0].PriceCurrent)
	assetRepo.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}

func TestGetBenchmark_SemHistoricoPortfolio_FallbackAveragePrice(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	portfolioQuotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 42.00}}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000.00},
		{Ticker: "SPX", CurrentValue: 5500.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)
	// portfolio sem histórico → fallback average_price
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(0.0, sql.ErrNoRows)
	priceHistory.On("GetPriceAtDate", "^BVSP", mock.AnythingOfType("time.Time")).Return(125000.00, nil)
	priceHistory.On("GetPriceAtDate", "SPX", mock.AnythingOfType("time.Time")).Return(5000.00, nil)

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.NoError(t, err)
	// portfolio: (4200-3800)/3800*100 ≈ 10.52%
	assert.InDelta(t, 10.52, result.PortfolioReturn, 0.01)
	assetRepo.AssertExpectations(t)
}

func TestGetBenchmark_PeriodoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "yearly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestGetBenchmark_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetBenchmark_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetBenchmark_BenchmarkProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	portfolioQuotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 42.00}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(nil, errors.New("api error"))
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)

	usecase := NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	result, err := usecase.Execute(userID, "weekly")

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}
