package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
)

type DeleteAssetUseCase struct {
	assetRepo domain.AssetRepository
}

func NewDeleteAssetUseCase(assetRepo domain.AssetRepository) DeleteAssetUseCase {
	return DeleteAssetUseCase{assetRepo: assetRepo}
}

func (d DeleteAssetUseCase) Execute(Ticker string) error {
	if Ticker == "" {
		return errors.New("Ticker doesnt not blank")
	}
	return d.assetRepo.DeleteByTicker(Ticker)
}
