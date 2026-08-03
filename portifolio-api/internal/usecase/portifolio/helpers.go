package portifolio

import (
	"portifolio-api/internal/domain"
	"sort"
)

func groupKey(groupBy string, asset domain.Asset) string {
	if groupBy == "sector" {
		if asset.Sector == "" {
			return "Outros"
		}
		return asset.Sector
	}
	return asset.Market
}

func buildMarketSummaries(assets []domain.Asset, priceOf func(ticker string) float64, groupBy string) []domain.MarketSummary {
	marketMap := make(map[string][]domain.AssetSummary)
	for _, asset := range assets {
		totalDay := priceOf(asset.Ticker) * float64(asset.Quantity)
		assetSummary := domain.AssetSummary{
			Ticker:   asset.Ticker,
			TotalDay: totalDay,
		}
		key := groupKey(groupBy, asset)
		marketMap[key] = append(marketMap[key], assetSummary)
	}

	marketSummaries := make([]domain.MarketSummary, 0, len(marketMap))
	for market, assetSummaries := range marketMap {
		var totalGainperMarket, totalLoserperMarket, totalperMarket float64
		for _, a := range assetSummaries {
			totalperMarket += a.TotalDay
			if a.TotalDay >= 0 {
				totalGainperMarket += a.TotalDay
			} else {
				totalLoserperMarket += a.TotalDay
			}
		}

		// top gainers — só positivos
		var gainers []domain.AssetSummary
		for _, a := range assetSummaries {
			if a.TotalDay >= 0 {
				gainers = append(gainers, a)
			}
		}
		sort.Slice(gainers, func(i, j int) bool {
			return gainers[i].TotalDay > gainers[j].TotalDay
		})
		topN := 3
		if len(gainers) < topN {
			topN = len(gainers)
		}
		topGainers := gainers[:topN]

		// top losers — só negativos
		var losers []domain.AssetSummary
		for _, a := range assetSummaries {
			if a.TotalDay < 0 {
				losers = append(losers, a)
			}
		}
		sort.Slice(losers, func(i, j int) bool {
			return losers[i].TotalDay < losers[j].TotalDay
		})
		topNL := 3
		if len(losers) < topNL {
			topNL = len(losers)
		}
		topLosers := losers[:topNL]

		marketSummaries = append(marketSummaries, domain.MarketSummary{
			MarketProvider:      market,
			TotalperMarket:      totalperMarket,
			TotalGainperMarket:  totalGainperMarket,
			TotalLoserperMarket: totalLoserperMarket,
			TopGainers:          topGainers,
			TopLosers:           topLosers,
		})
	}

	return marketSummaries
}
