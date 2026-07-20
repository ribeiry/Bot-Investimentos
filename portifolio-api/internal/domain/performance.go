package domain

type AssetPerformance struct {
	Ticker           string
	Market           string
	Quantity         float64
	AveragePrice     float64
	CurrentPrice     float64
	InvestedValue    float64
	CurrentValue     float64
	ProfitLoss       float64
	ReturnPercentage float64
}
