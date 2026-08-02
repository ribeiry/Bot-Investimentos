package domain

import "time"

type Asset struct {
	Ticker       string    `json:"ticker"`
	Market       string    `json:"market"`
	Quantity     float64   `json:"quantity"`
	AveragePrice float64   `json:"average_price"`
	Sector       string    `json:"sector,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}
