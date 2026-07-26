package portifolio

import (
	"errors"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAsset_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("DeleteByTicker", userID, "BBSE3").Return(nil)

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute(userID, "BBSE3")

	assert.NoError(t, err)
	assetRepo.AssertExpectations(t)
}

func TestDeleteAsset_TickerVazio(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute(userID, "")

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "DeleteByTicker")
}

func TestDeleteAsset_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	const userID int64 = 1

	assetRepo.On("DeleteByTicker", userID, "BBSE3").Return(errors.New("db error"))

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute(userID, "BBSE3")

	assert.Error(t, err)
	assetRepo.AssertExpectations(t)
}
