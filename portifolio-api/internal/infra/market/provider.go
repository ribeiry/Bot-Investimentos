package market

import (
	"log"
	"portifolio-api/internal/domain"
	"sync"
	"time"
)

type marketProviderWithFallback struct {
	brapi        domain.MarketProvider
	twelveData   domain.MarketProvider
	priceHistory domain.PriceHistoryRepository
}

func NewMarketProviderWithFallback(brapi domain.MarketProvider, twelveData domain.MarketProvider, priceHistory domain.PriceHistoryRepository) domain.MarketProvider {
	return marketProviderWithFallback{
		brapi:        brapi,
		twelveData:   twelveData,
		priceHistory: priceHistory,
	}
}

type providerResult struct {
	quotes []domain.Quote
	err    error
}

func (m marketProviderWithFallback) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var b3Assets []domain.Asset
	var usaAssets []domain.Asset

	log.Printf("[Provider] Total assets: %d", len(assets))
	for _, asset := range assets {
		log.Printf("[Provider] asset: %s market: '%s'", asset.Ticker, asset.Market)
		if asset.Market == "B3" {
			b3Assets = append(b3Assets, asset)
		} else {
			usaAssets = append(usaAssets, asset)
		}
	}
	log.Printf("[Provider] B3: %d, USA: %d", len(b3Assets), len(usaAssets))

	var wg sync.WaitGroup
	results := make(chan providerResult, 2)

	if len(b3Assets) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			quotes, err := m.brapi.GetByTickers(b3Assets)
			results <- providerResult{quotes: quotes, err: err}
		}()
	}

	if len(usaAssets) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			quotes, err := m.twelveData.GetByTickers(usaAssets)
			results <- providerResult{quotes: quotes, err: err}
		}()
	}

	wg.Wait()
	close(results)

	var quotes []domain.Quote
	for res := range results {
		if res.err != nil {
			return nil, res.err
		}
		quotes = append(quotes, res.quotes...)
	}

	for _, quote := range quotes {
		m.priceHistory.Save(domain.PriceHistory{
			Ticker:     quote.Ticker,
			Price:      quote.CurrentValue,
			CapturedAt: time.Now(),
		})
	}

	log.Printf("[Provider] Total quotes retornadas: %d", len(quotes))
	return quotes, nil
}
