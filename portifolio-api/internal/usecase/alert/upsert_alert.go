package alert

import (
	"errors"
	"portifolio-api/internal/domain"
)

type UpsertAlertUseCase struct {
	alertRepo domain.AlertRepository
}

func NewUpsertAlertUseCase(alertRepo domain.AlertRepository) UpsertAlertUseCase {
	return UpsertAlertUseCase{alertRepo: alertRepo}
}

type UpsertAlertInput struct {
	Ticker   string   `json:"ticker"`
	Market   string   `json:"market"`
	StopGain *float64 `json:"stop_gain"`
	StopLoss *float64 `json:"stop_loss"`
}

func (u UpsertAlertUseCase) Execute(userID int64, input UpsertAlertInput) error {
	if err := u.validate(input); err != nil {
		return err
	}
	return u.alertRepo.Upsert(domain.Alert{
		UserID:   userID,
		Ticker:   input.Ticker,
		Market:   input.Market,
		StopGain: input.StopGain,
		StopLoss: input.StopLoss,
	})
}

func (u UpsertAlertUseCase) validate(input UpsertAlertInput) error {
	if input.Ticker == "" {
		return errors.New("ticker is required")
	}
	if input.Market != "B3" && input.Market != "NASDAQ" && input.Market != "NYSE" {
		return errors.New("market must be B3, NASDAQ or NYSE")
	}
	if input.StopGain == nil && input.StopLoss == nil {
		return errors.New("at least one of stop_gain or stop_loss is required")
	}
	if input.StopGain != nil && *input.StopGain <= 0 {
		return errors.New("stop_gain must be greater than zero")
	}
	if input.StopLoss != nil && *input.StopLoss <= 0 {
		return errors.New("stop_loss must be greater than zero")
	}
	return nil
}
