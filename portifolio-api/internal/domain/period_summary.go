package domain

import "time"

type AssetPeriodSummary struct {
	Ticker         string  `json:"ticker"`
	Market         string  `json:"market"`
	Quantity       float64 `json:"quantity"`
	PriceStart     float64 `json:"price_start"`
	PriceCurrent   float64 `json:"price_current"`
	ValueStart     float64 `json:"value_start"`
	ValueCurrent   float64 `json:"value_current"`
	ChangeValue    float64 `json:"change_value"`
	ChangePercent  float64 `json:"change_percent"`
}

type PeriodSummary struct {
	Period             string               `json:"period"`
	PeriodStart        time.Time            `json:"period_start"`
	Assets             []AssetPeriodSummary `json:"assets"`
	TotalValueStart    float64              `json:"total_value_start"`
	TotalValueCurrent  float64              `json:"total_value_current"`
	TotalChangeValue   float64              `json:"total_change_value"`
	TotalChangePercent float64              `json:"total_change_percent"`
}
