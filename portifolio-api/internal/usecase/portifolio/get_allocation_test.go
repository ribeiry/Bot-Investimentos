package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllocation_ComSetor_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00, Sector: "Financeiro"},
		{Ticker: "ITSA4", Market: "B3", Quantity: 200, AveragePrice: 10.00, Sector: "Financeiro"},
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00, Sector: "Technology"},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 40.00},
		{Ticker: "ITSA4", CurrentValue: 12.00},
		{Ticker: "AAPL", CurrentValue: 160.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	// setor já preenchido → não chama UpdateSector

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	total := 40.0*100 + 12.0*200 + 160.0*10 // 4000 + 2400 + 1600 = 8000
	for _, alloc := range result {
		if alloc.Sector == "Financeiro" {
			assert.InDelta(t, 6400.0, alloc.TotalValue, 0.01)
			assert.InDelta(t, 80.0, alloc.Percentage, 0.01)
			assert.ElementsMatch(t, []string{"BBSE3", "ITSA4"}, alloc.Tickers)
		}
		if alloc.Sector == "Technology" {
			assert.InDelta(t, 1600.0, alloc.TotalValue, 0.01)
			assert.InDelta(t, 20.0, alloc.Percentage, 0.01)
		}
	}
	_ = total
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetAllocation_SemSetor_Outros(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "PETR4", Market: "B3", Quantity: 100, AveragePrice: 38.00, Sector: ""},
	}
	quotes := []domain.Quote{
		{Ticker: "PETR4", CurrentValue: 43.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Outros", result[0].Sector)
	assetRepo.AssertExpectations(t)
}

func TestGetAllocation_CarteiraVazia(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return([]domain.Asset{}, nil)

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	marketProvider.AssertNotCalled(t, "GetByTickers")
	assetRepo.AssertExpectations(t)
}

func TestGetAllocation_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assetRepo.On("ReturnAllPortfolio", userID).Return(nil, errors.New("db error"))

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
}

func TestGetAllocation_MarketProviderError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00, Sector: "Financeiro"},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestGetAllocation_PercentagemCorreta(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	assets := []domain.Asset{
		{Ticker: "A", Market: "B3", Quantity: 100, AveragePrice: 10.00, Sector: "Tech"},
		{Ticker: "B", Market: "B3", Quantity: 100, AveragePrice: 10.00, Sector: "Financeiro"},
		{Ticker: "C", Market: "B3", Quantity: 100, AveragePrice: 10.00, Sector: "Financeiro"},
	}
	quotes := []domain.Quote{
		{Ticker: "A", CurrentValue: 10.00},
		{Ticker: "B", CurrentValue: 10.00},
		{Ticker: "C", CurrentValue: 10.00},
	}

	assetRepo.On("ReturnAllPortfolio", userID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewGetAllocationUseCase(assetRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	for _, alloc := range result {
		if alloc.Sector == "Tech" {
			assert.InDelta(t, 33.33, alloc.Percentage, 0.01)
		}
		if alloc.Sector == "Financeiro" {
			assert.InDelta(t, 66.66, alloc.Percentage, 0.01)
		}
	}
	assetRepo.AssertExpectations(t)
}
