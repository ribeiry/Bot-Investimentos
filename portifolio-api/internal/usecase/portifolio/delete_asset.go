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

func (d DeleteAssetUseCase) Execute(userID int64, ticker string) error {
	if ticker == "" {
		return errors.New("Ticker doesnt not blank")
	}
	return d.assetRepo.DeleteByTicker(userID, ticker)
}
