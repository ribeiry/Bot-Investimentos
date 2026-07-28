package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	usecasemarket "portifolio-api/internal/usecase/market"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMarketRouter(h *MarketHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser)
	r.GET("/market/prices", h.GetPriceMarket)
	r.GET("/market/close", h.GetCloseMarket)
	return r
}

func buildMarketHandler(assetRepo *mocks.AssetRepository, marketProvider *mocks.MarketProvider) *MarketHandler {
	getClose := usecasemarket.NewGetCloseUseCase(assetRepo, marketProvider)
	getPrice := usecasemarket.NewGetPricesUseCase(assetRepo, marketProvider)
	return NewMarketHandler(getClose, getPrice)
}

// ─── GetPriceMarket ──────────────────────────────────────────────────────────

func TestMarketHandler_GetPrices_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupMarketRouter(buildMarketHandler(assetRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/market/prices", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "BBSE3")
	assetRepo.AssertExpectations(t)
}

func TestMarketHandler_GetPrices_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupMarketRouter(buildMarketHandler(assetRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/market/prices", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}

// ─── GetCloseMarket ──────────────────────────────────────────────────────────

func TestMarketHandler_GetClose_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupMarketRouter(buildMarketHandler(assetRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/market/close?market=B3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "BBSE3")
	assetRepo.AssertExpectations(t)
}

func TestMarketHandler_GetClose_MercadoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	r := setupMarketRouter(buildMarketHandler(assetRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/market/close?market=INVALID", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestMarketHandler_GetClose_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupMarketRouter(buildMarketHandler(assetRepo, marketProvider))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/market/close?market=B3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}
