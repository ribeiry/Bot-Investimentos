package alert

import (
	"errors"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAlert_Success(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	alertRepo.On("DeleteByTicker", userID, "BBSE3").Return(nil)

	usecase := NewDeleteAlertUseCase(alertRepo)
	err := usecase.Execute(userID, "BBSE3")

	assert.NoError(t, err)
	alertRepo.AssertExpectations(t)
}

func TestDeleteAlert_TickerVazio(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	usecase := NewDeleteAlertUseCase(alertRepo)
	err := usecase.Execute(userID, "")

	assert.Error(t, err)
	alertRepo.AssertNotCalled(t, "DeleteByTicker")
}

func TestDeleteAlert_RepoError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	alertRepo.On("DeleteByTicker", userID, "BBSE3").Return(errors.New("db error"))

	usecase := NewDeleteAlertUseCase(alertRepo)
	err := usecase.Execute(userID, "BBSE3")

	assert.Error(t, err)
	alertRepo.AssertExpectations(t)
}
