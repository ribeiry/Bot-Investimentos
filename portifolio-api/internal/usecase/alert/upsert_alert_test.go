package alert

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func floatPtr(v float64) *float64 { return &v }

func TestUpsertAlert_Success_StopGain(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}
	expected := domain.Alert{UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}

	alertRepo.On("Upsert", expected).Return(nil)

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.NoError(t, err)
	alertRepo.AssertExpectations(t)
}

func TestUpsertAlert_Success_StopLoss(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "AAPL", Market: "NYSE", StopLoss: floatPtr(140.00)}
	expected := domain.Alert{UserID: userID, Ticker: "AAPL", Market: "NYSE", StopLoss: floatPtr(140.00)}

	alertRepo.On("Upsert", expected).Return(nil)

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.NoError(t, err)
	alertRepo.AssertExpectations(t)
}

func TestUpsertAlert_Success_Ambos(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(50.00), StopLoss: floatPtr(30.00)}
	expected := domain.Alert{UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(50.00), StopLoss: floatPtr(30.00)}

	alertRepo.On("Upsert", expected).Return(nil)

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.NoError(t, err)
	alertRepo.AssertExpectations(t)
}

func TestUpsertAlert_TickerVazio(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "", Market: "B3", StopGain: floatPtr(45.00)}

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.Error(t, err)
	alertRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAlert_MercadoInvalido(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "INVALID", StopGain: floatPtr(45.00)}

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.Error(t, err)
	alertRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAlert_SemThreshold(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "B3"}

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.Error(t, err)
	alertRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAlert_StopGainZero(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(0)}

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.Error(t, err)
	alertRepo.AssertNotCalled(t, "Upsert")
}

func TestUpsertAlert_RepoError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	const userID int64 = 1

	input := UpsertAlertInput{Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}
	expected := domain.Alert{UserID: userID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}

	alertRepo.On("Upsert", expected).Return(errors.New("db error"))

	usecase := NewUpsertAlertUseCase(alertRepo)
	err := usecase.Execute(userID, input)

	assert.Error(t, err)
	alertRepo.AssertExpectations(t)
}
