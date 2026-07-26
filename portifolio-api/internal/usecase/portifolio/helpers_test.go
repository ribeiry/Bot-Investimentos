package portifolio

import (
	"portifolio-api/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildMarketSummaries_B3(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
		{Ticker: "ITSA4", Market: "B3", Quantity: 200, AveragePrice: 10.00},
	}

	priceOf := func(ticker string) float64 {
		prices := map[string]float64{
			"BBSE3": 40.00,
			"ITSA4": 12.00,
		}
		return prices[ticker]
	}

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result, 1)
	assert.Equal(t, "B3", result[0].MarketProvider)
	assert.InDelta(t, 6400.00, result[0].TotalperMarket, 0.01)
	assert.InDelta(t, 6400.00, result[0].TotalGainperMarket, 0.01)
	assert.Equal(t, 0.0, result[0].TotalLoserperMarket)
	assert.Len(t, result[0].TopGainers, 2)
}

func TestBuildMarketSummaries_MultiMarket(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50},
		{Ticker: "AAPL", Market: "NYSE", Quantity: 10, AveragePrice: 150.00},
	}

	priceOf := func(ticker string) float64 {
		prices := map[string]float64{
			"BBSE3": 40.00,
			"AAPL":  160.00,
		}
		return prices[ticker]
	}

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result, 2)
}

func TestBuildMarketSummaries_DoisAtivos(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 40.00},
		{Ticker: "ITSA4", Market: "B3", Quantity: 200, AveragePrice: 12.00},
	}

	priceOf := func(ticker string) float64 {
		prices := map[string]float64{
			"BBSE3": 35.00,
			"ITSA4": 14.00,
		}
		return prices[ticker]
	}

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result, 1)
	assert.InDelta(t, 6300.00, result[0].TotalperMarket, 0.01)
	assert.InDelta(t, 6300.00, result[0].TotalGainperMarket, 0.01)
	assert.Equal(t, 0.0, result[0].TotalLoserperMarket)
	assert.Len(t, result[0].TopGainers, 2)
	assert.Len(t, result[0].TopLosers, 0)
}

func TestBuildMarketSummaries_TopNLimitado(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "A1", Market: "B3", Quantity: 1, AveragePrice: 1.0},
		{Ticker: "A2", Market: "B3", Quantity: 2, AveragePrice: 1.0},
		{Ticker: "A3", Market: "B3", Quantity: 3, AveragePrice: 1.0},
		{Ticker: "A4", Market: "B3", Quantity: 4, AveragePrice: 1.0},
	}

	priceOf := func(ticker string) float64 { return 10.0 }

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result[0].TopGainers, 3)
}
func TestBuildMarketSummaries_ComPerdas(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 40.00},
		{Ticker: "ITSA4", Market: "B3", Quantity: 200, AveragePrice: 12.00},
	}

	priceOf := func(ticker string) float64 {
		prices := map[string]float64{
			"BBSE3": -5.00, // negativo para forçar o caminho de losers
			"ITSA4": 14.00,
		}
		return prices[ticker]
	}

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result, 1)
	assert.Len(t, result[0].TopLosers, 1)
	assert.Equal(t, "BBSE3", result[0].TopLosers[0].Ticker)
	assert.Len(t, result[0].TopGainers, 1)
	assert.Equal(t, "ITSA4", result[0].TopGainers[0].Ticker)
}
func TestBuildMarketSummaries_MultiplasPercas(t *testing.T) {
	assets := []domain.Asset{
		{Ticker: "A1", Market: "B3", Quantity: 100, AveragePrice: 40.00},
		{Ticker: "A2", Market: "B3", Quantity: 200, AveragePrice: 12.00},
		{Ticker: "A3", Market: "B3", Quantity: 50, AveragePrice: 10.00},
	}

	priceOf := func(ticker string) float64 {
		prices := map[string]float64{
			"A1": -3.00,
			"A2": -2.00,
			"A3": -1.00,
		}
		return prices[ticker]
	}

	result := buildMarketSummaries(assets, priceOf)

	assert.Len(t, result[0].TopLosers, 3)
}
