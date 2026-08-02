package portifolio

import (
	"database/sql"
	"portifolio-api/internal/domain"
)

var benchmarkAssets = []struct {
	Name   string
	Ticker string
	Market string
}{
	{Name: "IBOV", Ticker: "^BVSP", Market: "B3"},
	{Name: "S&P500", Ticker: "SPX", Market: "NYSE"},
}

type GetBenchmarkUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
	priceHistory   domain.PriceHistoryRepository
}

func NewGetBenchmarkUseCase(
	assetRepo domain.AssetRepository,
	marketProvider domain.MarketProvider,
	priceHistory domain.PriceHistoryRepository,
) GetBenchmarkUseCase {
	return GetBenchmarkUseCase{
		assetRepo:      assetRepo,
		marketProvider: marketProvider,
		priceHistory:   priceHistory,
	}
}

func (g GetBenchmarkUseCase) Execute(userID int64, period string) (*domain.BenchmarkComparison, error) {
	periodStart, err := resolvePeriodStart(period)
	if err != nil {
		return nil, err
	}

	// ── Portfolio return ─────────────────────────────────────────────────────
	assets, err := g.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}

	portfolioQuotes, err := g.marketProvider.GetByTickers(assets)
	if err != nil {
		return nil, err
	}
	currentPriceMap := make(map[string]float64, len(portfolioQuotes))
	for _, q := range portfolioQuotes {
		currentPriceMap[q.Ticker] = q.CurrentValue
	}

	var portfolioValueStart, portfolioValueCurrent float64
	for _, a := range assets {
		priceStart, err := g.priceHistory.GetPriceAtDate(a.Ticker, periodStart)
		if err != nil {
			if err == sql.ErrNoRows {
				priceStart = a.AveragePrice
			} else {
				return nil, err
			}
		}
		portfolioValueStart += a.Quantity * priceStart
		portfolioValueCurrent += a.Quantity * currentPriceMap[a.Ticker]
	}

	var portfolioReturn float64
	if portfolioValueStart > 0 {
		portfolioReturn = ((portfolioValueCurrent - portfolioValueStart) / portfolioValueStart) * 100
	}

	// ── Benchmark returns ────────────────────────────────────────────────────
	benchAssets := make([]domain.Asset, len(benchmarkAssets))
	for i, b := range benchmarkAssets {
		benchAssets[i] = domain.Asset{Ticker: b.Ticker, Market: b.Market}
	}

	benchQuotes, err := g.marketProvider.GetByTickers(benchAssets)
	if err != nil {
		return nil, err
	}
	benchCurrentMap := make(map[string]float64, len(benchQuotes))
	for _, q := range benchQuotes {
		benchCurrentMap[q.Ticker] = q.CurrentValue
	}

	var results []domain.BenchmarkResult
	for _, b := range benchmarkAssets {
		priceCurrent := benchCurrentMap[b.Ticker]

		priceStart, err := g.priceHistory.GetPriceAtDate(b.Ticker, periodStart)
		if err != nil {
			if err == sql.ErrNoRows {
				results = append(results, domain.BenchmarkResult{
					Name:         b.Name,
					Ticker:       b.Ticker,
					PriceCurrent: priceCurrent,
					NoHistory:    true,
				})
				continue
			}
			return nil, err
		}

		var returnPct float64
		if priceStart > 0 {
			returnPct = ((priceCurrent - priceStart) / priceStart) * 100
		}

		results = append(results, domain.BenchmarkResult{
			Name:                b.Name,
			Ticker:              b.Ticker,
			PriceStart:          priceStart,
			PriceCurrent:        priceCurrent,
			ReturnPercent:       returnPct,
			RelativePerformance: portfolioReturn - returnPct,
		})
	}

	return &domain.BenchmarkComparison{
		Period:          period,
		PeriodStart:     periodStart,
		PortfolioReturn: portfolioReturn,
		Benchmarks:      results,
	}, nil
}
