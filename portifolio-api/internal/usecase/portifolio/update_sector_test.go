package portifolio

import (
	"errors"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateSector_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("UpdateSector", userID, "PETR4", "Energia").Return(nil)

	usecase := NewUpdateSectorUseCase(assetRepo)
	err := usecase.Execute(userID, "PETR4", UpdateSectorInput{Sector: "Energia"})

	assert.NoError(t, err)
	assetRepo.AssertExpectations(t)
}

func TestUpdateSector_Limpar(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("UpdateSector", userID, "PETR4", "").Return(nil)

	usecase := NewUpdateSectorUseCase(assetRepo)
	err := usecase.Execute(userID, "PETR4", UpdateSectorInput{Sector: ""})

	assert.NoError(t, err)
	assetRepo.AssertExpectations(t)
}

func TestUpdateSector_TickerVazio(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	usecase := NewUpdateSectorUseCase(assetRepo)
	err := usecase.Execute(userID, "", UpdateSectorInput{Sector: "Energia"})

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "UpdateSector")
}

func TestUpdateSector_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("UpdateSector", userID, "PETR4", "Energia").Return(errors.New("db error"))

	usecase := NewUpdateSectorUseCase(assetRepo)
	err := usecase.Execute(userID, "PETR4", UpdateSectorInput{Sector: "Energia"})

	assert.Error(t, err)
	assetRepo.AssertExpectations(t)
}
