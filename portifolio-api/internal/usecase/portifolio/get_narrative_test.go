package portifolio

import (
	"context"
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testPromptTemplate = "prompt with {{PORTFOLIO_DATA}}"

func buildNarrativeUseCase(
	assetRepo *mocks.AssetRepository,
	provider *mocks.MarketProvider,
	priceHistory *mocks.PriceHistoryRepository,
	llm *mocks.LLMProvider,
	cache *mocks.NarrativeCache,
	userLimiter *mocks.RateLimiter,
	globalLimiter *mocks.GlobalRateLimiter,
	deadline time.Duration,
) GetNarrativeUseCase {
	return NewGetNarrativeUseCase(
		NewGetSummaryUseCase(assetRepo, provider, provider),
		NewGetPerformanceUseCase(assetRepo, provider),
		NewGetBenchmarkUseCase(assetRepo, provider, priceHistory),
		NewGetAllocationUseCase(assetRepo, provider),
		llm,
		cache,
		userLimiter,
		globalLimiter,
		testPromptTemplate,
		deadline,
	)
}

func stubAggregate(assetRepo *mocks.AssetRepository, provider *mocks.MarketProvider, priceHistory *mocks.PriceHistoryRepository) {
	assets := []domain.Asset{
		{Ticker: "BBSE3", Market: "B3", Quantity: 100, AveragePrice: 30, Sector: "Financeiro"},
	}
	quotes := []domain.Quote{
		{Ticker: "BBSE3", CurrentValue: 35},
	}
	benchAssets := []domain.Asset{
		{Ticker: "^BVSP", Market: "B3"},
		{Ticker: "SPX", Market: "NYSE"},
	}
	benchQuotes := []domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000},
		{Ticker: "SPX", CurrentValue: 5500},
	}

	// summary + performance + benchmark + allocation all call ReturnAllPortfolio
	assetRepo.On("ReturnAllPortfolio", int64(1)).Return(assets, nil)
	provider.On("GetByTickers", assets).Return(quotes, nil)
	provider.On("GetByTickers", benchAssets).Return(benchQuotes, nil)
	priceHistory.On("GetPriceAtDate", mock.Anything, mock.Anything).Return(float64(28), nil)
}

func TestGetNarrative_CacheHit(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("resumo cacheado", true)

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 500*time.Millisecond)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "resumo cacheado", text)
	llm.AssertNotCalled(t, "GenerateNarrative")
	assetRepo.AssertNotCalled(t, "ReturnAllPortfolio")
}

func TestGetNarrative_UserRateLimit(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(false)

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 500*time.Millisecond)

	text, err := uc.Execute(context.Background(), 1)

	assert.ErrorIs(t, err, domain.ErrRateLimitExceeded)
	assert.Empty(t, text)
	cache.AssertNotCalled(t, "Get")
	llm.AssertNotCalled(t, "GenerateNarrative")
}

func TestGetNarrative_CarteiraVazia_NaoChamaLLM(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)

	assetRepo.On("ReturnAllPortfolio", int64(1)).Return([]domain.Asset{}, nil)
	provider.On("GetByTickers", []domain.Asset{}).Return([]domain.Quote{}, nil)
	provider.On("GetByTickers", []domain.Asset{
		{Ticker: "^BVSP", Market: "B3"},
		{Ticker: "SPX", Market: "NYSE"},
	}).Return([]domain.Quote{
		{Ticker: "^BVSP", CurrentValue: 130000},
		{Ticker: "SPX", CurrentValue: 5500},
	}, nil)
	priceHistory.On("GetPriceAtDate", mock.Anything, mock.Anything).Return(float64(0), nil)

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 500*time.Millisecond)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, emptyPortfolioText, text)
	llm.AssertNotCalled(t, "GenerateNarrative")
	globalLimiter.AssertNotCalled(t, "Allow")
}

func TestGetNarrative_GlobalRateLimit_FallbackSilent(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	stubAggregate(assetRepo, provider, priceHistory)
	globalLimiter.On("Allow").Return(false)

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 500*time.Millisecond)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, fallbackText, text)
	llm.AssertNotCalled(t, "GenerateNarrative")
}

func TestGetNarrative_LLMRapido_ValidoECacheia(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	stubAggregate(assetRepo, provider, priceHistory)
	globalLimiter.On("Allow").Return(true)

	narrative := "📊 Sua carteira **subiu 2%**. BBSE3 puxou a alta."
	llm.On("GenerateNarrative", mock.Anything, mock.Anything).Return(narrative, nil)
	cache.On("Set", int64(1), narrative, narrativeCacheTTL).Return()

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 2*time.Second)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, narrative, text)
	// dá tempo para a goroutine chamar Set caso ainda esteja em curso
	time.Sleep(50 * time.Millisecond)
	cache.AssertCalled(t, "Set", int64(1), narrative, narrativeCacheTTL)
}

func TestGetNarrative_LLMComTickerInvalido_Fallback_NaoCacheia(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	stubAggregate(assetRepo, provider, priceHistory)
	globalLimiter.On("Allow").Return(true)

	narrative := "XYZW9 disparou hoje" // ticker fabricado
	llm.On("GenerateNarrative", mock.Anything, mock.Anything).Return(narrative, nil)

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 2*time.Second)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, fallbackText, text)
	time.Sleep(50 * time.Millisecond)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetNarrative_LLMErro_Fallback(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	stubAggregate(assetRepo, provider, priceHistory)
	globalLimiter.On("Allow").Return(true)
	llm.On("GenerateNarrative", mock.Anything, mock.Anything).Return("", errors.New("groq down"))

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 2*time.Second)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, fallbackText, text)
	time.Sleep(50 * time.Millisecond)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetNarrative_LLMLento_RetornaProcessing_CacheiaDepois(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	stubAggregate(assetRepo, provider, priceHistory)
	globalLimiter.On("Allow").Return(true)

	narrative := "📊 BBSE3 subiu 2%"
	llm.On("GenerateNarrative", mock.Anything, mock.Anything).Run(func(mock.Arguments) {
		time.Sleep(150 * time.Millisecond)
	}).Return(narrative, nil)
	cache.On("Set", int64(1), narrative, narrativeCacheTTL).Return()

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 50*time.Millisecond)

	start := time.Now()
	text, err := uc.Execute(context.Background(), 1)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, processingText, text)
	assert.Less(t, elapsed, 130*time.Millisecond, "deve retornar antes do LLM terminar")

	// aguarda a goroutine terminar e gravar no cache
	time.Sleep(200 * time.Millisecond)
	cache.AssertCalled(t, "Set", int64(1), narrative, narrativeCacheTTL)
}

func TestGetNarrative_AggregateErro_Fallback(t *testing.T) {
	assetRepo := new(mocks.AssetRepository)
	provider := new(mocks.MarketProvider)
	priceHistory := new(mocks.PriceHistoryRepository)
	llm := new(mocks.LLMProvider)
	cache := new(mocks.NarrativeCache)
	userLimiter := new(mocks.RateLimiter)
	globalLimiter := new(mocks.GlobalRateLimiter)

	userLimiter.On("Allow", int64(1)).Return(true)
	cache.On("Get", int64(1)).Return("", false)
	assetRepo.On("ReturnAllPortfolio", int64(1)).Return(nil, errors.New("db down"))
	provider.On("GetByTickers", mock.Anything).Return([]domain.Quote{}, nil).Maybe()
	priceHistory.On("GetPriceAtDate", mock.Anything, mock.Anything).Return(float64(0), nil).Maybe()

	uc := buildNarrativeUseCase(assetRepo, provider, priceHistory, llm, cache, userLimiter, globalLimiter, 500*time.Millisecond)

	text, err := uc.Execute(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, fallbackText, text)
	llm.AssertNotCalled(t, "GenerateNarrative")
}
