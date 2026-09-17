package portifolio

import (
	"portifolio-api/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeTotals_SomaInvestedECurrent(t *testing.T) {
	performance := []domain.AssetPerformance{
		{Ticker: "AAPL", InvestedValue: 1800, CurrentValue: 3370},
		{Ticker: "BBSE3", InvestedValue: 3000, CurrentValue: 4041},
		{Ticker: "MSFT", InvestedValue: 2000, CurrentValue: 2488.75},
		{Ticker: "PETR4", InvestedValue: 1750, CurrentValue: 2430.5},
	}

	totals := computeTotals(performance)

	assert.InDelta(t, 8550, totals.Invested, 0.001)
	assert.InDelta(t, 12330.25, totals.Current, 0.001)
	assert.InDelta(t, 3780.25, totals.ProfitLoss, 0.001)
	assert.InDelta(t, 44.213, totals.ReturnPercentage, 0.01)
}

func TestComputeTotals_PerformanceVazio_TotaisZerados(t *testing.T) {
	totals := computeTotals(nil)

	assert.Equal(t, float64(0), totals.Invested)
	assert.Equal(t, float64(0), totals.Current)
	assert.Equal(t, float64(0), totals.ProfitLoss)
	assert.Equal(t, float64(0), totals.ReturnPercentage)
}

func TestComputeTotals_InvestedZero_NaoDivideePorZero(t *testing.T) {
	performance := []domain.AssetPerformance{
		{Ticker: "X", InvestedValue: 0, CurrentValue: 100},
	}

	totals := computeTotals(performance)

	assert.Equal(t, float64(0), totals.Invested)
	assert.Equal(t, float64(100), totals.Current)
	assert.Equal(t, float64(100), totals.ProfitLoss)
	assert.Equal(t, float64(0), totals.ReturnPercentage)
}

func TestComputeTotals_PrejuizoNegativo(t *testing.T) {
	performance := []domain.AssetPerformance{
		{Ticker: "X", InvestedValue: 1000, CurrentValue: 800},
	}

	totals := computeTotals(performance)

	assert.InDelta(t, -200, totals.ProfitLoss, 0.001)
	assert.InDelta(t, -20, totals.ReturnPercentage, 0.001)
}
