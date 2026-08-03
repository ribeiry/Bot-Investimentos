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

func (g GetSummaryUseCase) Execute(userID int64, mode string, groupBy string) (*domain.PortfolioSummary, error) {
	assets, err := g.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}

	provider := g.cachedProvider
	if mode == "realtime" {
		provider = g.marketProvider
	}

	quotes, err := provider.GetByTickers(assets)
	if err != nil {
		return nil, err
	}

	quoteMap := make(map[string]domain.Quote)
	for _, quote := range quotes {
		quoteMap[quote.Ticker] = quote
	}

	if groupBy == "" {
		groupBy = "market"
	}
	summaries := buildMarketSummaries(assets, func(ticker string) float64 {
		return quoteMap[ticker].CurrentValue
	}, groupBy)

	return &domain.PortfolioSummary{
		MarketSummary: summaries,
		QuotedAt:      time.Now(),
	}, nil
}
