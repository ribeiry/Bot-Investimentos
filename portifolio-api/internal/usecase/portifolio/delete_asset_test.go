package portifolio

import (
	"errors"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAsset_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	assetRepo.On("DeleteByTicker", "BBSE3").Return(nil)

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute("BBSE3")

	assert.NoError(t, err)
	assetRepo.AssertExpectations(t)
}

func TestDeleteAsset_TickerVazio(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute("")

	assert.Error(t, err)
	assetRepo.AssertNotCalled(t, "DeleteByTicker")
}

func TestDeleteAsset_RepoError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)

	assetRepo.On("DeleteByTicker", "BBSE3").Return(errors.New("db error"))

	usecase := NewDeleteAssetUseCase(assetRepo)
	err := usecase.Execute("BBSE3")

	assert.Error(t, err)
	assetRepo.AssertExpectations(t)
}
