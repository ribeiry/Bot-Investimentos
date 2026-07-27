package alert

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAlerts_StopGainDisparado(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
	}
	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3"}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 46.00}}

	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "BBSE3", result[0].Ticker)
	assert.Equal(t, "STOP_GAIN", result[0].TriggerType)
	assert.Equal(t, 46.00, result[0].CurrentPrice)
	alertRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestCheckAlerts_StopLossDisparado(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "AAPL", Market: "NYSE", StopLoss: floatPtr(140.00), Active: true},
	}
	assets := []domain.Asset{{Ticker: "AAPL", Market: "NYSE"}}
	quotes := []domain.Quote{{Ticker: "AAPL", CurrentValue: 135.00}}

	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "AAPL", result[0].Ticker)
	assert.Equal(t, "STOP_LOSS", result[0].TriggerType)
	assert.Equal(t, 135.00, result[0].CurrentPrice)
	alertRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestCheckAlerts_NenhumDisparado(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(50.00), StopLoss: floatPtr(30.00), Active: true},
	}
	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3"}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}

	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	alertRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestCheckAlerts_MultiplosTickers(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
		{ID: 2, UserID: userID, Ticker: "AAPL", Market: "NYSE", StopLoss: floatPtr(140.00), Active: true},
		{ID: 3, UserID: userID, Ticker: "ITSA4", Market: "B3", StopGain: floatPtr(15.00), Active: true},
	}
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3"},
		{Ticker: "AAPL", Market: "NYSE"},
		{Ticker: "ITSA4", Market: "B3"},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 46.00}, // dispara STOP_GAIN
		{Ticker: "AAPL", CurrentValue: 155.00}, // não dispara
		{Ticker: "ITSA4", CurrentValue: 12.00}, // não dispara
	}

	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "BBSE3", result[0].Ticker)
	assert.Equal(t, "STOP_GAIN", result[0].TriggerType)
	alertRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}

func TestCheckAlerts_SemAlertas(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alertRepo.On("GetAllByUserID", userID).Return([]domain.Alert{}, nil)

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 0)
	marketProvider.AssertNotCalled(t, "GetByTickers")
	alertRepo.AssertExpectations(t)
}

func TestCheckAlerts_RepoError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alertRepo.On("GetAllByUserID", userID).Return(nil, errors.New("db error"))

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	alertRepo.AssertExpectations(t)
}

func TestCheckAlerts_MarketProviderError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)
	const userID int64 = 1

	alerts := []domain.Alert{
		{ID: 1, UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
	}
	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3"}}

	alertRepo.On("GetAllByUserID", userID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(nil, errors.New("api error"))

	usecase := NewCheckAlertsUseCase(alertRepo, marketProvider)
	result, err := usecase.Execute(userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	alertRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
}
