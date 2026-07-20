package domain

import "time"

type PortfolioSummary struct {
	MarketSummary []MarketSummary
	QuotedAt      time.Time
}

type MarketSummary struct {
	MarketProvider      string
	TotalperMarket      float64
	TotalLoserperMarket float64
	TotalGainperMarket  float64
	TopLosers           []AssetSummary
	TopGainers          []AssetSummary
}

type AssetSummary struct {
	Ticker   string
	TotalDay float64
}
