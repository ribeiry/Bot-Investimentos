package market

import "portifolio-api/internal/domain"

type GetPricesUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
}

func NewGetPricesUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider) GetPricesUseCase {
	return GetPricesUseCase{assetRepo: assetRepo, marketProvider: marketProvider}
}

func (g GetPricesUseCase) Execute() ([]domain.Quote, error) {

	assetsRepo, err := g.assetRepo.ReturnAllPortfolio()

	if err != nil {
		return nil, err
	}

	return g.marketProvider.GetByTickers(assetsRepo)

}
