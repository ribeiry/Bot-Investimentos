package portifolio

import (
	"database/sql"
	"errors"
	"portifolio-api/internal/domain"
	"time"
)

type GetPeriodSummaryUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
	priceHistory   domain.PriceHistoryRepository
}

func NewGetPeriodSummaryUseCase(
	assetRepo domain.AssetRepository,
	marketProvider domain.MarketProvider,
	priceHistory domain.PriceHistoryRepository,
) GetPeriodSummaryUseCase {
	return GetPeriodSummaryUseCase{
		assetRepo:      assetRepo,
		marketProvider: marketProvider,
		priceHistory:   priceHistory,
	}
}

func (g GetPeriodSummaryUseCase) Execute(userID int64, period string) (*domain.PeriodSummary, error) {
	periodStart, err := resolvePeriodStart(period)
	if err != nil {
		return nil, err
	}

	assets, err := g.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}

	quotes, err := g.marketProvider.GetByTickers(assets)
	if err != nil {
		return nil, err
	}
	currentPriceMap := make(map[string]float64, len(quotes))
	for _, q := range quotes {
		currentPriceMap[q.Ticker] = q.CurrentValue
	}

	var assetSummaries []domain.AssetPeriodSummary
	var totalStart, totalCurrent float64

	for _, asset := range assets {
		priceStart, err := g.priceHistory.GetPriceAtDate(asset.Ticker, periodStart)
		if err != nil {
			if err == sql.ErrNoRows {
				priceStart = asset.AveragePrice
			} else {
				return nil, err
			}
		}

		priceCurrent := currentPriceMap[asset.Ticker]
		valueStart := asset.Quantity * priceStart
		valueCurrent := asset.Quantity * priceCurrent
		changeValue := valueCurrent - valueStart

		var changePercent float64
		if valueStart != 0 {
			changePercent = (changeValue / valueStart) * 100
		}

		totalStart += valueStart
		totalCurrent += valueCurrent

		assetSummaries = append(assetSummaries, domain.AssetPeriodSummary{
			Ticker:        asset.Ticker,
			Market:        asset.Market,
			Quantity:      asset.Quantity,
			PriceStart:    priceStart,
			PriceCurrent:  priceCurrent,
			ValueStart:    valueStart,
			ValueCurrent:  valueCurrent,
			ChangeValue:   changeValue,
			ChangePercent: changePercent,
		})
	}

	totalChangeValue := totalCurrent - totalStart
	var totalChangePercent float64
	if totalStart != 0 {
		totalChangePercent = (totalChangeValue / totalStart) * 100
	}

	return &domain.PeriodSummary{
		Period:             period,
		PeriodStart:        periodStart,
		Assets:             assetSummaries,
		TotalValueStart:    totalStart,
		TotalValueCurrent:  totalCurrent,
		TotalChangeValue:   totalChangeValue,
		TotalChangePercent: totalChangePercent,
	}, nil
}

func resolvePeriodStart(period string) (time.Time, error) {
	now := time.Now()
	switch period {
	case "weekly":
		// volta até a segunda-feira da semana atual
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		monday := now.AddDate(0, 0, -(weekday - 1))
		return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location()), nil
	case "monthly":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()), nil
	default:
		return time.Time{}, errors.New("period must be 'weekly' or 'monthly'")
	}
}
