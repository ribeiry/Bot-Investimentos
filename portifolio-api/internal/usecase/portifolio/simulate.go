package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
)

type SimulateUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
	priceHistory   domain.PriceHistoryRepository
}

func NewSimulateUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider, priceHistory domain.PriceHistoryRepository) SimulateUseCase {
	return SimulateUseCase{assetRepo: assetRepo, marketProvider: marketProvider, priceHistory: priceHistory}
}

type SimulateInput struct {
	Operations []domain.SimulationOperation `json:"operations"`
}

func (s SimulateUseCase) Execute(userID int64, input SimulateInput) (*domain.SimulationResult, error) {
	if err := s.validate(input); err != nil {
		return nil, err
	}

	assets, err := s.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}

	// build ticker list: portfolio + buy operations (single API call)
	tickerSet := make(map[string]domain.Asset)
	for _, a := range assets {
		tickerSet[a.Ticker] = a
	}
	for _, op := range input.Operations {
		if op.Action == "buy" {
			if _, exists := tickerSet[op.Ticker]; !exists {
				tickerSet[op.Ticker] = domain.Asset{Ticker: op.Ticker, Market: op.Market}
			}
		}
	}
	allAssets := make([]domain.Asset, 0, len(tickerSet))
	for _, a := range tickerSet {
		allAssets = append(allAssets, a)
	}

	quoteMap := s.buildQuoteMap(allAssets)

	// snapshot atual
	current := buildSnapshot(assets, quoteMap)

	// copia em memória para simular
	simAssets := make(map[string]domain.Asset, len(assets))
	for _, a := range assets {
		simAssets[a.Ticker] = a
	}

	var skipped []domain.SimulationSkipped

	for _, op := range input.Operations {
		switch op.Action {
		case "sell":
			existing, exists := simAssets[op.Ticker]
			if !exists {
				skipped = append(skipped, domain.SimulationSkipped{
					Ticker: op.Ticker, Action: "sell", Reason: "asset not in portfolio",
				})
				continue
			}
			if op.Quantity > existing.Quantity {
				skipped = append(skipped, domain.SimulationSkipped{
					Ticker: op.Ticker, Action: "sell", Reason: "insufficient quantity",
				})
				continue
			}
			existing.Quantity -= op.Quantity
			if existing.Quantity == 0 {
				delete(simAssets, op.Ticker)
			} else {
				simAssets[op.Ticker] = existing
			}

		case "buy":
			price, hasQuote := quoteMap[op.Ticker]
			if !hasQuote || price == 0 {
				skipped = append(skipped, domain.SimulationSkipped{
					Ticker: op.Ticker, Action: "buy", Reason: "quote unavailable",
				})
				continue
			}
			if existing, exists := simAssets[op.Ticker]; exists {
				// média ponderada do preço médio
				totalQty := existing.Quantity + op.Quantity
				existing.AveragePrice = (existing.Quantity*existing.AveragePrice + op.Quantity*price) / totalQty
				existing.Quantity = totalQty
				simAssets[op.Ticker] = existing
			} else {
				simAssets[op.Ticker] = domain.Asset{
					Ticker:       op.Ticker,
					Market:       op.Market,
					Quantity:     op.Quantity,
					AveragePrice: price,
				}
			}
		}
	}

	simAssetList := make([]domain.Asset, 0, len(simAssets))
	for _, a := range simAssets {
		simAssetList = append(simAssetList, a)
	}
	simulated := buildSnapshot(simAssetList, quoteMap)

	if skipped == nil {
		skipped = []domain.SimulationSkipped{}
	}

	return &domain.SimulationResult{
		Skipped:   skipped,
		Current:   current,
		Simulated: simulated,
		Delta: domain.SimulationDelta{
			Value:         simulated.TotalCurrentValue - current.TotalCurrentValue,
			ReturnPercent: simulated.TotalReturnPercent - current.TotalReturnPercent,
		},
	}, nil
}

func (s SimulateUseCase) validate(input SimulateInput) error {
	if len(input.Operations) == 0 {
		return errors.New("operations list cannot be empty")
	}
	for _, op := range input.Operations {
		if op.Action != "buy" && op.Action != "sell" {
			return errors.New("action must be 'buy' or 'sell'")
		}
		if op.Ticker == "" {
			return errors.New("ticker is required")
		}
		if op.Quantity <= 0 {
			return errors.New("quantity must be greater than zero")
		}
		if op.Action == "buy" && op.Market == "" {
			return errors.New("market is required for buy operations")
		}
	}
	return nil
}

// buildQuoteMap: price_history primeiro; chama API só para o que faltar
func (s SimulateUseCase) buildQuoteMap(assets []domain.Asset) map[string]float64 {
	quoteMap := make(map[string]float64, len(assets))
	var missing []domain.Asset

	for _, a := range assets {
		price, err := s.priceHistory.GetLastPrice(a.Ticker)
		if err == nil && price > 0 {
			quoteMap[a.Ticker] = price
		} else {
			missing = append(missing, a)
		}
	}

	if len(missing) > 0 {
		quotes, err := s.marketProvider.GetByTickers(missing)
		if err == nil {
			for _, q := range quotes {
				quoteMap[q.Ticker] = q.CurrentValue
			}
		}
	}

	return quoteMap
}

func buildSnapshot(assets []domain.Asset, quoteMap map[string]float64) domain.SimulationSnapshot {
	var invested, current float64
	for _, a := range assets {
		invested += a.Quantity * a.AveragePrice
		current += a.Quantity * quoteMap[a.Ticker]
	}
	var returnPct float64
	if invested > 0 {
		returnPct = ((current - invested) / invested) * 100
	}
	return domain.SimulationSnapshot{
		TotalInvested:      invested,
		TotalCurrentValue:  current,
		TotalReturnPercent: returnPct,
	}
}
