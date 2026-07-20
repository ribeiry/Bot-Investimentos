package portifolio

import (
	"portifolio-api/internal/domain"
)

type GetAssetUseCase struct {
	assetRepo domain.AssetRepository
}

func NewGetAssetUseCase(assetRepo domain.AssetRepository) GetAssetUseCase {
	return GetAssetUseCase{assetRepo: assetRepo}
}

func (g GetAssetUseCase) Execute() ([]domain.Asset, error) {

	return g.assetRepo.ReturnAllPortfolio()
}
