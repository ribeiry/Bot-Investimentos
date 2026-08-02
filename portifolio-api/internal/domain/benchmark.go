package domain

import "time"

type BenchmarkResult struct {
	Name                string  `json:"name"`
	Ticker              string  `json:"ticker"`
	PriceStart          float64 `json:"price_start"`
	PriceCurrent        float64 `json:"price_current"`
	ReturnPercent       float64 `json:"return_percent"`
	RelativePerformance float64 `json:"relative_performance"` // portfolio_return - benchmark_return
	NoHistory           bool    `json:"no_history,omitempty"`
}

type BenchmarkComparison struct {
	Period          string            `json:"period"`
	PeriodStart     time.Time         `json:"period_start"`
	PortfolioReturn float64           `json:"portfolio_return_percent"`
	Benchmarks      []BenchmarkResult `json:"benchmarks"`
}
