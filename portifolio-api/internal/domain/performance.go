package domain

type AssetPerformance struct {
	Ticker           string  `json:"ticker"`
	Market           string  `json:"market"`
	Sector           string  `json:"sector,omitempty"`
	Quantity         float64 `json:"quantity"`
	AveragePrice     float64 `json:"average_price"`
	CurrentPrice     float64 `json:"current_price"`
	InvestedValue    float64 `json:"invested_value"`
	CurrentValue     float64 `json:"current_value"`
	ProfitLoss       float64 `json:"profit_loss"`
	ReturnPercentage float64 `json:"return_percentage"`
}
