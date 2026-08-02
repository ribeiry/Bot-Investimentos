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
	"portifolio-api/internal/usecase/portifolio"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupPortfolioRouter(h *PortfolioHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser)
	r.GET("/portfolio/assets", h.GetAssets)
	r.POST("/portfolio/assets", h.UpsertAsset)
	r.DELETE("/portfolio/assets/:ticker", h.DeleteAsset)
	r.GET("/portfolio/summary", h.GetSummaryAsset)
	r.GET("/portfolio/performance", h.GetPerformance)
	r.GET("/portfolio/period-summary", h.GetPeriodSummary)
	r.GET("/portfolio/benchmark", h.GetBenchmark)
	r.GET("/portfolio/allocation", h.GetAllocation)
	return r
}

func buildPortfolioHandler(assetRepo *mocks.AssetRepository, marketProvider *mocks.MarketProvider, priceHistory *mocks.PriceHistoryRepository) *PortfolioHandler {
	upsert := portifolio.NewUpsertAssetUseCase(assetRepo)
	delete := portifolio.NewDeleteAssetUseCase(assetRepo)
	getAsset := portifolio.NewGetAssetUseCase(assetRepo)
	getSummary := portifolio.NewGetSummaryUseCase(assetRepo, marketProvider, marketProvider)
	getPerformance := portifolio.NewGetPerformanceUseCase(assetRepo, marketProvider)
	getPeriodSummary := portifolio.NewGetPeriodSummaryUseCase(assetRepo, marketProvider, priceHistory)
	getBenchmark := portifolio.NewGetBenchmarkUseCase(assetRepo, marketProvider, priceHistory)
	getAllocation := portifolio.NewGetAllocationUseCase(assetRepo, marketProvider)
	return NewPortfolioHandler(upsert, delete, getAsset, getSummary, getPerformance, getPeriodSummary, getBenchmark, getAllocation)
}

// ─── GetAssets ───────────────────────────────────────────────────────────────

func TestPortfolioHandler_GetAssets_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/assets", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "BBSE3")
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_GetAssets_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/assets", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "error")
	assetRepo.AssertExpectations(t)
}

// ─── UpsertAsset ─────────────────────────────────────────────────────────────

func TestPortfolioHandler_UpsertAsset_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	asset := domain.Asset{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}
	assetRepo.On("Upsert", testUserID, asset).Return(nil)

	body, _ := json.Marshal(asset)
	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/portfolio/assets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "message")
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_UpsertAsset_BadJSON(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/portfolio/assets", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
}

func TestPortfolioHandler_UpsertAsset_ValidationError(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	asset := domain.Asset{Ticker: "", Market: "B3", Quantity: 100, AveragePrice: 38.50}
	body, _ := json.Marshal(asset)
	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/portfolio/assets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertNotCalled(t, "Upsert")
}

// ─── DeleteAsset ─────────────────────────────────────────────────────────────

func TestPortfolioHandler_DeleteAsset_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("DeleteByTicker", testUserID, "BBSE3").Return(nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/portfolio/assets/BBSE3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "message")
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_DeleteAsset_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("DeleteByTicker", testUserID, "BBSE3").Return(errors.New("db error"))

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/portfolio/assets/BBSE3", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}

// ─── GetSummary ──────────────────────────────────────────────────────────────

func TestPortfolioHandler_GetSummary_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/summary?mode=realtime", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "data")
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_GetSummary_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/summary?mode=realtime", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}

// ─── GetPerformance ──────────────────────────────────────────────────────────

func TestPortfolioHandler_GetPerformance_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.50}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 40.00}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/performance", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "BBSE3")
	assetRepo.AssertExpectations(t)
}

// ─── GetPeriodSummary ────────────────────────────────────────────────────────

func TestPortfolioHandler_GetPeriodSummary_Weekly_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00}}
	quotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 42.00}}
	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/period-summary?period=weekly", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "weekly")
	assert.Contains(t, w.Body.String(), "BBSE3")
	assetRepo.AssertExpectations(t)
	marketProvider.AssertExpectations(t)
	priceHistory.AssertExpectations(t)
}

func TestPortfolioHandler_GetPeriodSummary_PeriodoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/period-summary?period=yearly", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

// ─── GetBenchmark ────────────────────────────────────────────────────────────

func TestPortfolioHandler_GetBenchmark_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	assets := []domain.Asset{{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00}}
	portfolioQuotes := []domain.Quote{{Ticker: "BBSE3", CurrentValue: 42.00}}
	benchAssets := []domain.Asset{
		{Ticker: "^BVSP", Market: "B3"},
		{Ticker: "SPX", Market: "NYSE"},
	}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000.00},
		{Ticker: "SPX", CurrentValue: 5500.00},
	}

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(portfolioQuotes, nil)
	marketProvider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)
	priceHistory.On("GetPriceAtDate", "BBSE3", mock.AnythingOfType("time.Time")).Return(40.00, nil)
	priceHistory.On("GetPriceAtDate", "^BVSP", mock.AnythingOfType("time.Time")).Return(125000.00, nil)
	priceHistory.On("GetPriceAtDate", "SPX", mock.AnythingOfType("time.Time")).Return(5000.00, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/benchmark?period=weekly", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "portfolio_return_percent")
	assert.Contains(t, w.Body.String(), "IBOV")
	assert.Contains(t, w.Body.String(), "SPX")
	assetRepo.AssertExpectations(t)
}

// ─── GetAllocation ───────────────────────────────────────────────────────────

func TestPortfolioHandler_GetAllocation_Success(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 38.00, Sector: "Financeiro"},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 40.00, Sector: "Financeiro"},
	}

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(assets, nil)
	marketProvider.On("GetByTickers", assets).Return(quotes, nil)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/allocation", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assert.Contains(t, w.Body.String(), "Financeiro")
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_GetAllocation_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/allocation", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}

func TestPortfolioHandler_GetBenchmark_PeriodoInvalido(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, priceHistory))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/benchmark?period=yearly", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestPortfolioHandler_GetPerformance_Error(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	marketProvider := new(mocks.MarketProvider)

	assetRepo.On("ReturnAllPortfolio", testUserID).Return(nil, errors.New("db error"))

	r := setupPortfolioRouter(buildPortfolioHandler(assetRepo, marketProvider, new(mocks.PriceHistoryRepository)))
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/performance", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), testTelegramID)
	assetRepo.AssertExpectations(t)
}
