package portifolio

import (
	"portifolio-api/internal/domain"
	"time"
)

type GetSummaryUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
	cachedProvider domain.MarketProvider
}

func NewGetSummaryUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider, cachedProvider domain.MarketProvider) GetSummaryUseCase {
	return GetSummaryUseCase{
		assetRepo:      assetRepo,
		marketProvider: marketProvider,
		cachedProvider: cachedProvider,
	}
}
func (g GetSummaryUseCase) Execute(mode string) (*domain.PortfolioSummary, error) {

	assets, err := g.assetRepo.ReturnAllPortfolio()

	if err != nil {
		return nil, err
	}
	if mode == "realtime" {

		quotes, err := g.marketProvider.GetByTickers(assets)

		if err != nil {
			return nil, err
		}
		quoteMap := make(map[string]domain.Quote)
		////monta a lista de tickers para consultar o Yahoo Finance
		for _, quote := range quotes {
			quoteMap[quote.Ticker] = quote
		}
		summaries := buildMarketSummaries(assets, func(ticker string) float64 {
			return quoteMap[ticker].CurrentValue
		})

		summary := &domain.PortfolioSummary{
			MarketSummary: summaries,
			QuotedAt:      time.Now(),
		}

		return summary, nil
	} else {
		quotes, err := g.cachedProvider.GetByTickers(assets)
		if err != nil {
			return nil, err
		}

		quoteMap := make(map[string]domain.Quote)
		for _, quote := range quotes {
			quoteMap[quote.Ticker] = quote
		}

		summaries := buildMarketSummaries(assets, func(ticker string) float64 {
			return quoteMap[ticker].CurrentValue
		})

		summary := &domain.PortfolioSummary{
			MarketSummary: summaries,
			QuotedAt:      time.Now(),
		}
		return summary, nil
	}
}
