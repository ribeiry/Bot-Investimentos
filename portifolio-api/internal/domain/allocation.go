package domain

type SectorAllocation struct {
	Sector     string   `json:"sector"`
	TotalValue float64  `json:"total_value"`
	Percentage float64  `json:"percentage"`
	Tickers    []string `json:"tickers"`
}
