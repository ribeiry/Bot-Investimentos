package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	usecasealert "portifolio-api/internal/usecase/alert"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAlertRouter(h *AlertHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser)
	r.GET("/alerts", h.GetAlerts)
	r.POST("/alerts", h.UpsertAlert)
	r.DELETE("/alerts/:ticker", h.DeleteAlert)
	r.GET("/alerts/check", h.CheckAlerts)
	return r
}

func buildAlertHandler(alertRepo *mocks.AlertRepository, marketProvider *mocks.MarketProvider) *AlertHandler {
	upsert := usecasealert.NewUpsertAlertUseCase(alertRepo)
	delete := usecasealert.NewDeleteAlertUseCase(alertRepo)
	get := usecasealert.NewGetAlertsUseCase(alertRepo)
	check := usecasealert.NewCheckAlertsUseCase(alertRepo, marketProvider)
	return NewAlertHandler(upsert, delete, get, check)
}

// ─── GetAlerts ───────────────────────────────────────────────────────────────

func TestAlertHandler_GetAlerts_Success(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alerts := []domain.Alert{
		{ID: 1, UserID: testUserID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
	}
	alertRepo.On("GetAllByUserID", testUserID).Return(alerts, nil)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/alerts", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "BBSE3")
	alertRepo.AssertExpectations(t)
}

func TestAlertHandler_GetAlerts_Error(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alertRepo.On("GetAllByUserID", testUserID).Return(nil, errors.New("db error"))

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/alerts", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	alertRepo.AssertExpectations(t)
}

// ─── UpsertAlert ─────────────────────────────────────────────────────────────

func TestAlertHandler_UpsertAlert_Success(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	expected := domain.Alert{UserID: testUserID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}
	alertRepo.On("Upsert", expected).Return(nil)

	input := usecasealert.UpsertAlertInput{Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00)}
	body, _ := json.Marshal(input)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/alerts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "message")
	alertRepo.AssertExpectations(t)
}

func TestAlertHandler_UpsertAlert_BadJSON(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/alerts", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	alertRepo.AssertNotCalled(t, "Upsert")
}

func TestAlertHandler_UpsertAlert_ValidationError(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	input := usecasealert.UpsertAlertInput{Ticker: "BBSE3", Market: "B3"} // sem stop_gain nem stop_loss
	body, _ := json.Marshal(input)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/alerts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	alertRepo.AssertNotCalled(t, "Upsert")
}

// ─── DeleteAlert ─────────────────────────────────────────────────────────────

func TestAlertHandler_DeleteAlert_Success(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alertRepo.On("DeleteByTicker", testUserID, "BBSE3").Return(nil)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/alerts/BBSE3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "message")
	alertRepo.AssertExpectations(t)
}

func TestAlertHandler_DeleteAlert_Error(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alertRepo.On("DeleteByTicker", testUserID, "BBSE3").Return(errors.New("db error"))

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/alerts/BBSE3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	alertRepo.AssertExpectations(t)
}

// ─── CheckAlerts ─────────────────────────────────────────────────────────────

func TestAlertHandler_CheckAlerts_ComTrigger(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alerts := []domain.Alert{
		{ID: 1, UserID: testUserID, Ticker: "BBSE3", Market: "B3", StopGain: floatPtr(45.00), Active: true},
	}
	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3"}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 46.00}}

	alertRepo.On("GetAllByUserID", testUserID).Return(alerts, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/alerts/check", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "STOP_GAIN")
	alertRepo.AssertExpectations(t)
}

func TestAlertHandler_CheckAlerts_SemTrigger(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alertRepo.On("GetAllByUserID", testUserID).Return([]domain.Alert{}, nil)

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/alerts/check", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "data")
	alertRepo.AssertExpectations(t)
}

func TestAlertHandler_CheckAlerts_Error(t *testing.T) {
	alertRepo := new(mocks.AlertRepository)
	marketProvider := new(mocks.MarketProvider)

	alertRepo.On("GetAllByUserID", testUserID).Return(nil, errors.New("db error"))

	r := setupAlertRouter(buildAlertHandler(alertRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/alerts/check", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	alertRepo.AssertExpectations(t)
}
