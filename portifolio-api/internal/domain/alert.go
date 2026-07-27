package domain

import "time"

type Alert struct {
	ID        int64
	UserID    int64
	Ticker    string
	Market    string
	StopGain  *float64
	StopLoss  *float64
	Active    bool
	CreatedAt time.Time
}

type TriggeredAlert struct {
	Alert
	CurrentPrice float64
	TriggerType  string // "STOP_GAIN" or "STOP_LOSS"
}

type AlertRepository interface {
	Upsert(alert Alert) error
	GetAllByUserID(userID int64) ([]Alert, error)
	DeleteByTicker(userID int64, ticker string) error
}
