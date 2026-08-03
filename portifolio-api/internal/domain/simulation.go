package domain

type SimulationOperation struct {
	Action   string  `json:"action"`
	Ticker   string  `json:"ticker"`
	Market   string  `json:"market"`
	Quantity float64 `json:"quantity"`
}

type SimulationSkipped struct {
	Ticker string `json:"ticker"`
	Action string `json:"action"`
	Reason string `json:"reason"`
}

type SimulationSnapshot struct {
	TotalInvested      float64 `json:"total_invested"`
	TotalCurrentValue  float64 `json:"total_current_value"`
	TotalReturnPercent float64 `json:"total_return_percent"`
}

type SimulationDelta struct {
	Value         float64 `json:"value"`
	ReturnPercent float64 `json:"return_percent"`
}

type SimulationResult struct {
	Skipped   []SimulationSkipped `json:"skipped"`
	Current   SimulationSnapshot  `json:"current"`
	Simulated SimulationSnapshot  `json:"simulated"`
	Delta     SimulationDelta     `json:"delta"`
}
