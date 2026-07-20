package portifolio

import "portifolio-api/internal/domain"

type GetPerformanceUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
}

func NewGetPerformanceUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider) GetPerformanceUseCase {
	return GetPerformanceUseCase{assetRepo: assetRepo, marketProvider: marketProvider}
}

func (g GetPerformanceUseCase) Execute() ([]domain.AssetPerformance, error) {

	assets, error := g.assetRepo.ReturnAllPortfolio()

	if error != nil {
		return nil, error
	}

	quotes, error := g.marketProvider.GetByTickers(assets)
	if error != nil {
		return nil, error
	}
	quoteMap := make(map[string]domain.Quote)
	for _, quote := range quotes {
		quoteMap[quote.Ticker] = quote

	}

	var performances []domain.AssetPerformance
	for _, asset := range assets {

		valorInvestido := asset.Quantity * asset.AveragePrice
		valorAtual := asset.Quantity * quoteMap[asset.Ticker].CurrentValue
		valorDiferencaAtual := valorAtual - valorInvestido
		performanceTicker := (valorDiferencaAtual / valorInvestido) * 100

		performances = append(performances, domain.AssetPerformance{
			Ticker:           asset.Ticker,
			Market:           asset.Market,
			Quantity:         asset.Quantity,
			AveragePrice:     asset.AveragePrice,
			CurrentPrice:     quoteMap[asset.Ticker].CurrentValue,
			InvestedValue:    valorInvestido,
			CurrentValue:     valorAtual,
			ProfitLoss:       valorDiferencaAtual,
			ReturnPercentage: performanceTicker,
		})
	}

	return performances, nil

}
