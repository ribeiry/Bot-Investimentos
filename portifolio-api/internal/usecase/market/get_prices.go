package market

import "portifolio-api/internal/domain"

type GetPricesUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
}

func NewGetPricesUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider) GetPricesUseCase {
	return GetPricesUseCase{assetRepo: assetRepo, marketProvider: marketProvider}
}

func (g GetPricesUseCase) Execute(userID int64) ([]domain.Quote, error) {
	assets, err := g.assetRepo.ReturnAllPortfolio(userID)
	if err != nil {
		return nil, err
	}
	return g.marketProvider.GetByTickers(assets)
}
