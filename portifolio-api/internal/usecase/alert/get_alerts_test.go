package alert

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAlerts_Success(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
		{ID: 2, UserID: userID, Ticker: "AAPL", Market: "NYSE", StopLoss: floatPtr(140.00), Active: true},
	}
	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)

	usecase := NewGetAlertsUseCase(alertRepo)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "BBSE3", result[0].Ticker)
	alertRepo.AssertExpectations(t)
}

func TestGetAlerts_ListaVazia(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	alertRepo.On("GetAllByUserID", userID).Return([]domain.Alert{}, nil)

	usecase := NewGetAlertsUseCase(alertRepo)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	alertRepo.AssertExpectations(t)
}

func TestGetAlerts_RepoError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	alertRepo.On("GetAllByUserID", userID).Return(nil, errors.New("db error"))

	usecase := NewGetAlertsUseCase(alertRepo)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	alertRepo.AssertExpectations(t)
}
