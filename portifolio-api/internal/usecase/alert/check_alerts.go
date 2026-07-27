package alert

import "portifolio-api/internal/domain"

type CheckAlertsUseCase struct {
	alertRepo      domain.AlertRepository
	marketProvider domain.MarketProvider
}

func NewCheckAlertsUseCase(alertRepo domain.AlertRepository, marketProvider domain.MarketProvider) CheckAlertsUseCase {
	return CheckAlertsUseCase{alertRepo: alertRepo, marketProvider: marketProvider}
}

func (c CheckAlertsUseCase) Execute(userID int64) ([]domain.TriggeredAlert, error) {
	alerts, err := c.alertRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(alerts) == 0 {
		return []domain.TriggeredAlert{}, nil
	}

	assets := make([]domain.Asset, len(alerts))
	for i, a := range alerts {
		assets[i] = domain.Asset{Ticker: a.Ticker, Market: a.Market}
	}

	quotes, err := c.marketProvider.GetByTickers(assets)
	if err != nil {
		return nil, err
	}

	quoteMap := make(map[string]float64, len(quotes))
	for _, quote := range quotes {
		quoteMap[quote.Ticker] = quote.CurrentValue
	}

	var triggered []domain.TriggeredAlert
	for _, alert := range alerts {
		price := quoteMap[alert.Ticker]
		if alert.StopGain != nil && price >= *alert.StopGain {
			triggered = append(triggered, domain.TriggeredAlert{
				Alert:        alert,
				CurrentPrice: price,
				TriggerType:  "STOP_GAIN",
			})
		} else if alert.StopLoss != nil && price <= *alert.StopLoss {
			triggered = append(triggered, domain.TriggeredAlert{
				Alert:        alert,
				CurrentPrice: price,
				TriggerType:  "STOP_LOSS",
			})
		}
	}

	return triggered, nil
}
