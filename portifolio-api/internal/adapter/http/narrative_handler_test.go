package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portifolio-api/internal/mocks"
	"portifolio-api/internal/usecase/portifolio"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildNarrativeHandler(
	assetRepo *mocks.AssetRepository,
	provider *mocks.MarketProvider,
	priceHistory *mocks.PriceHistoryRepository,
	llm *mocks.LLMProvider,
	cache *mocks.NarrativeCache,
	userLimiter *mocks.RateLimiter,
	globalLimiter *mocks.GlobalRateLimiter,
) *PortfolioHandler {
	upsert := portifolio.NewUpsertAssetUseCase(assetRepo)
	del := portifolio.NewDeleteAssetUseCase(assetRepo)
	getAsset := portifolio.NewGetAssetUseCase(assetRepo)
	summary := portifolio.NewGetSummaryUseCase(assetRepo, provider, provider)
	perf := portifolio.NewGetPerformanceUseCase(assetRepo, provider)
	period := portifolio.NewGetPeriodSummaryUseCase(assetRepo, provider, priceHistory)
	bench := portifolio.NewGetBenchmarkUseCase(assetRepo, provider, priceHistory)
	alloc := portifolio.NewGetAllocationUseCase(assetRepo, provider)
	sector := portifolio.NewUpdateSectorUseCase(assetRepo)
	sim := portifolio.NewSimulateUseCase(assetRepo, provider, priceHistory)

	narrative := portifolio.NewGetNarrativeUseCase(
		summary, perf, bench, alloc,
		llm, cache, userLimiter, globalLimiter,
		"prompt with {{PORTFOLIO_DATA}}",
		2*time.Second,
	)
	return NewPortfolioHandler(upsert, del, getAsset, summary, perf, period, bench, alloc, sector, sim, narrative)
}

func setupNarrativeRouter(h *PortfolioHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(injectUser)
	r.GET("/portfolio/summary/narrative", h.GetNarrative)
	return r
}

func TestPortfolioHandler_GetNarrative_CacheHit(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", testUserID).Return(true)
	cache.On("Get", testUserID).Return("📊 resumo pronto", true)

	h := buildNarrativeHandler(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter)
	r := setupNarrativeRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/summary/narrative", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, testTelegramID, body["telegram_id"])
	data := body["data"].(map[string]any)
	assert.Equal(t, "📊 resumo pronto", data["text"])
}

func TestPortfolioHandler_GetNarrative_RateLimit429(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", testUserID).Return(false)

	h := buildNarrativeHandler(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter)
	r := setupNarrativeRouter(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/portfolio/summary/narrative", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	var body map[string]any
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, testTelegramID, body["telegram_id"])
	assert.Contains(t, body["error"], "limite")

	llm.AssertNotCalled(t, "GenerateNarrative")
	cache.AssertNotCalled(t, "Get", mock.Anything)
}
