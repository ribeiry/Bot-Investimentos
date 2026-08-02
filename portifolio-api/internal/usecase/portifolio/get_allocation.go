package portifolio

import "portifolio-api/internal/domain"

type GetAllocationUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
}

func NewGetAllocationUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider) GetAllocationUseCase {
	return GetAllocationUseCase{assetRepo: assetRepo, marketProvider: marketProvider}
}

func (g GetAllocationUseCase) Execute(userID int64) ([]domain.SectorAllocation, error) {
	assets, err := g.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		return []domain.SectorAllocation{}, nil
	}

	quotes, err := g.marketProvider.GetByTickers(assets)
	if err != nil {
		return nil, err
	}

	quoteMap := make(map[string]domain.Quote, len(quotes))
	for _, q := range quotes {
		quoteMap[q.Ticker] = q
	}

	for i, a := range assets {
		if a.Sector == "" {
			if q, ok := quoteMap[a.Ticker]; ok && q.Sector != "" {
				g.assetRepo.UpdateSector(userID, a.Ticker, q.Sector)
				assets[i].Sector = q.Sector
			}
		}
	}

	sectorMap := make(map[string]*domain.SectorAllocation)
	var totalValue float64

	for _, a := range assets {
		sector := a.Sector
		if sector == "" {
			sector = "Outros"
		}
		value := a.Quantity * quoteMap[a.Ticker].CurrentValue
		totalValue += value

		if _, ok := sectorMap[sector]; !ok {
			sectorMap[sector] = &domain.SectorAllocation{Sector: sector}
		}
		sectorMap[sector].TotalValue += value
		sectorMap[sector].Tickers = append(sectorMap[sector].Tickers, a.Ticker)
	}

	result := make([]domain.SectorAllocation, 0, len(sectorMap))
	for _, alloc := range sectorMap {
		if totalValue > 0 {
			alloc.Percentage = (alloc.TotalValue / totalValue) * 100
		}
		result = append(result, *alloc)
	}

	return result, nil
}
