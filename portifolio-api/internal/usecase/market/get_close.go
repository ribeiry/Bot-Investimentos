package market

import (
	"errors"
	"portifolio-api/internal/domain"
)

type GetCloseUseCase struct {
	assetRepo      domain.AssetRepository
	marketProvider domain.MarketProvider
}

func NewGetCloseUseCase(assetRepo domain.AssetRepository, marketProvider domain.MarketProvider) GetCloseUseCase {
	return GetCloseUseCase{assetRepo: assetRepo, marketProvider: marketProvider}
}

func (g GetCloseUseCase) Execute(market string) ([]domain.Quote, error) {

	if err := g.validate(market); err != nil {
		return nil, err
	}
	assetsRepo, err := g.assetRepo.ReturnAllPortfolio()

	if err != nil {
		return nil, err
	}

	return g.marketProvider.GetByTickers(assetsRepo)

}

func (g GetCloseUseCase) validate(market string) error {
	if market != "B3" && market != "NASDAQ" && market != "NYSE" {
		return errors.New("Market is wrong")
	}
	return nil
}
