package domain

import "time"

type Quote struct {
	MarketName    string
	Ticker        string
	QuotedAt      time.Time
	PreviousValue float64
	CurrentValue  float64
	Sector        string
}
