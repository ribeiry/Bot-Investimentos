package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
)

type UpsertAssetUseCase struct {
	assetRepo domain.AssetRepository
}

func NewUpsertAssetUseCase(assetRepo domain.AssetRepository) UpsertAssetUseCase {
	return UpsertAssetUseCase{assetRepo: assetRepo}
}

func (u UpsertAssetUseCase) Execute(asset domain.Asset) error {
	if err := u.validate(asset); err != nil {
		return err
	}
	return u.assetRepo.Upsert(asset)
}

func (u UpsertAssetUseCase) validate(asset domain.Asset) error {

	if asset.Ticker == "" {
		return errors.New("Error the Ticker is nill or Empty")
	}
	if asset.Quantity == 0 || asset.Quantity < 0 {
		return errors.New("Quantity doesnt 0 or < 0 ")
	}
	if asset.AveragePrice == 0 || asset.AveragePrice < 0 {
		return errors.New("AveragePrice doesnt 0 or < 0 ")
	}
	if asset.Market != "B3" && asset.Market != "NASDAQ" && asset.Market != "NYSE" {
		return errors.New("Market is wrong")
	}
	return nil
}
