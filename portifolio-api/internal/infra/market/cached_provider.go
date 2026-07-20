package market

import (
	"portifolio-api/internal/domain"
	"time"
)

type cachedMarketProvider struct {
	priceHistory domain.PriceHistoryRepository
}

func NewCachedMarketProvider(priceHistory domain.PriceHistoryRepository) domain.MarketProvider {
	return cachedMarketProvider{priceHistory: priceHistory}
}

func (c cachedMarketProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var quotes []domain.Quote
	for _, asset := range assets {
		price, err := c.priceHistory.GetLastPrice(asset.Ticker)
		if err != nil {
			continue
		}
		quotes = append(quotes, domain.Quote{
			Ticker:       asset.Ticker,
			CurrentValue: price,
			QuotedAt:     time.Now(),
		})
	}
	return quotes, nil
}
