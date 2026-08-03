package portifolio

import (
	"errors"
	"portifolio-api/internal/domain"
)

type UpdateSectorUseCase struct {
	assetRepo domain.AssetRepository
}

func NewUpdateSectorUseCase(assetRepo domain.AssetRepository) UpdateSectorUseCase {
	return UpdateSectorUseCase{assetRepo: assetRepo}
}

type UpdateSectorInput struct {
	Sector string `json:"sector"`
}

func (u UpdateSectorUseCase) Execute(userID int64, ticker string, input UpdateSectorInput) error {
	if ticker == "" {
		return errors.New("ticker is required")
	}
	return u.assetRepo.UpdateSector(userID, ticker, input.Sector)
}
