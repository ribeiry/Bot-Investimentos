package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// helper: cria use case com os três mocks
func newSimulateUC(assetRepo *mocks.AssetRepository, mp *mocks.MarketProvider, ph *mocks.PriceHistoryRepository) SimulateUseCase {
	return NewSimulateUseCase(assetRepo, mp, ph)
}

func TestSimulate_VendaECompra_PriceHistoryHit(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "BBSE3", Quantity: 50},
		{Action: "buy", Ticker: "AAPL", Market: "NYSE", Quantity: 5},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	// BBSE3 no price_history, AAPL não
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)
	ph.On("GetLastPrice", "AAPL").Return(0.0, errors.New("not found"))
	// API só para AAPL
	marketProvider.On("GetByTickers", mock.MatchedBy(func(a []domain.Asset) bool {
		return len(a) == 1 && a[0].Ticker == "AAPL"
	})).Return([]domain.Quote{
		{Ticker: "AAPL", CurrentValue: 160.00},
	}, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 0)
	assert.InDelta(t, 3800.00, result.Current.TotalInvested, 0.01)
	assert.InDelta(t, 4000.00, result.Current.TotalCurrentValue, 0.01)
	// 50 BBSE3 (40) + 5 AAPL (160) = 2000 + 800 = 2800
	assert.InDelta(t, 2800.00, result.Simulated.TotalCurrentValue, 0.01)

	assetRepo.AssertExpectations(t)
	ph.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestSimulate_VendaECompra_TodosNoHistorico(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "BBSE3", Quantity: 50},
		{Action: "buy", Ticker: "AAPL", Market: "NYSE", Quantity: 5},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	// ambos no price_history — API não deve ser chamada
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)
	ph.On("GetLastPrice", "AAPL").Return(160.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 0)
	assert.InDelta(t, 2800.00, result.Simulated.TotalCurrentValue, 0.01)

	marketProvider.AssertNotCalled(t, "GetByTickers")
	ph.AssertExpectations(t)
}

func TestSimulate_VendaInsuficiente_Skipped(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "BBSE3", Quantity: 200},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 1)
	assert.Equal(t, "BBSE3", result.Skipped[0].Ticker)
	assert.Equal(t, "insufficient quantity", result.Skipped[0].Reason)
	assert.Equal(t, result.Current.TotalCurrentValue, result.Simulated.TotalCurrentValue)

	marketProvider.AssertNotCalled(t, "GetByTickers")
}

func TestSimulate_AtivoForaCarteira_Skipped(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "XPTO3", Quantity: 10},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 1)
	assert.Equal(t, "XPTO3", result.Skipped[0].Ticker)
	assert.Equal(t, "asset not in portfolio", result.Skipped[0].Reason)

	marketProvider.AssertNotCalled(t, "GetByTickers")
}

func TestSimulate_CompraSemCotacao_Skipped(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "buy", Ticker: "XPTO3", Market: "B3", Quantity: 10},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)
	ph.On("GetLastPrice", "XPTO3").Return(0.0, errors.New("not found"))
	marketProvider.On("GetByTickers", mock.MatchedBy(func(a []domain.Asset) bool {
		return len(a) == 1 && a[0].Ticker == "XPTO3"
	})).Return([]domain.Quote{}, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 1)
	assert.Equal(t, "XPTO3", result.Skipped[0].Ticker)
	assert.Equal(t, "quote unavailable", result.Skipped[0].Reason)
}

func TestSimulate_CompraAtivoExistente_MediaPonderada(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "buy", Ticker: "BBSE3", Market: "B3", Quantity: 100},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.Len(t, result.Skipped, 0)
	// 200 BBSE3 * 40 = 8000
	assert.InDelta(t, 8000.00, result.Simulated.TotalCurrentValue, 0.01)
	// avg = (100*38 + 100*40) / 200 = 39
	assert.InDelta(t, 200*39.0, result.Simulated.TotalInvested, 0.01)

	marketProvider.AssertNotCalled(t, "GetByTickers")
}

func TestSimulate_VendaTodosOsAtivos(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "BBSE3", Quantity: 100},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	assert.InDelta(t, 0.0, result.Simulated.TotalCurrentValue, 0.01)
	assert.InDelta(t, 0.0, result.Simulated.TotalInvested, 0.01)

	marketProvider.AssertNotCalled(t, "GetByTickers")
}

func TestSimulate_OperacoesVazias_Erro(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, SimulateInput{
		Operations: []domain.SimulationOperation{},
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestSimulate_BuySemMarket_Erro(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "buy", Ticker: "AAPL", Quantity: 5},
	}}

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestSimulate_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, SimulateInput{
		Operations: []domain.SimulationOperation{
			{Action: "sell", Ticker: "BBSE3", Quantity: 10},
		},
	})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSimulate_Delta_Correto(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	ph := new(mocks.PriceHistoryRepository)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00},
	}
	input := SimulateInput{Operations: []domain.SimulationOperation{
		{Action: "sell", Ticker: "BBSE3", Quantity: 50},
	}}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	ph.On("GetLastPrice", "BBSE3").Return(40.00, nil)

	result, err := newSimulateUC(assetRepo, marketProvider, ph).Execute(userID, input)

	assert.NoError(t, err)
	// current: 4000, simulated: 2000 → delta = -2000
	assert.InDelta(t, -2000.00, result.Delta.Value, 0.01)

	marketProvider.AssertNotCalled(t, "GetByTickers")
}
