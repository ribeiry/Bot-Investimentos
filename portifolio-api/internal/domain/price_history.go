package domain

import "time"

type PriceHistory struct {
	Ticker     string
	Price      float64
	CapturedAt time.Time
}

type PriceHistoryRepository interface {
	Save(history PriceHistory) error
	GetLastPrice(ticker string) (float64, error)
	GetPriceAtDate(ticker string, date time.Time) (float64, error)
}
