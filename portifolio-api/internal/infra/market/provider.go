package market

import (
	"log"
	"portifolio-api/internal/domain"
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
	var quotes []domain.Quote

	if len(b3Assets) > 0 {
		b3Quotes, err := m.brapi.GetByTickers(b3Assets)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, b3Quotes...)
	}

	if len(usaAssets) > 0 {
		usaQuotes, err := m.twelveData.GetByTickers(usaAssets)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, usaQuotes...)
	}

	for _, quote := range quotes {
		m.priceHistory.Save(domain.PriceHistory{
			Ticker:     quote.Ticker,
			Price:      quote.CurrentValue,
			CapturedAt: time.Now(),
		})
	}
	// após as chamadas:
	log.Printf("[Provider] Total quotes retornadas: %d", len(quotes))

	return quotes, nil
}
