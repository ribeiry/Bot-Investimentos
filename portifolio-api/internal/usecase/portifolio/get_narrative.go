package portifolio

import (
	"context"
	"errors"
	"log"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/infra/cache"
	"portifolio-api/internal/infra/ratelimit"
	"time"
)

const (
	narrativeCacheTTL     = 1 * time.Hour
	narrativeLLMTimeout   = 10 * time.Second
	fallbackText          = "Não foi possível gerar o resumo agora. Tente novamente em alguns minutos."
	processingText        = "⏳ Estamos processando seu resumo. Tente novamente em alguns segundos."
	emptyPortfolioText    = "Sua carteira está vazia. Cadastre ativos para receber um resumo."
	promptDataPlaceholder = "{{PORTFOLIO_DATA}}"
)

var (
	errEmptyPortfolio = errors.New("empty portfolio")
	errGlobalLimit    = errors.New("global limit exceeded")
)

type narrativeOutcome struct {
	text    string
	tickers []string
	err     error
}

type GetNarrativeUseCase struct {
	getSummary     GetSummaryUseCase
	getPerformance GetPerformanceUseCase
	getBenchmark   GetBenchmarkUseCase
	getAllocation  GetAllocationUseCase
	llm            domain.LLMProvider
	cache          cache.NarrativeCache
	userLimiter    ratelimit.RateLimiter
	globalLimiter  ratelimit.GlobalRateLimiter
	promptTemplate string
	deadline       time.Duration
}

func NewGetNarrativeUseCase(
	summary GetSummaryUseCase,
	performance GetPerformanceUseCase,
	benchmark GetBenchmarkUseCase,
	allocation GetAllocationUseCase,
	llm domain.LLMProvider,
	narrativeCache cache.NarrativeCache,
	userLimiter ratelimit.RateLimiter,
	globalLimiter ratelimit.GlobalRateLimiter,
	promptTemplate string,
	deadline time.Duration,
) GetNarrativeUseCase {
	return GetNarrativeUseCase{
		getSummary:     summary,
		getPerformance: performance,
		getBenchmark:   benchmark,
		getAllocation:  allocation,
		llm:            llm,
		cache:          narrativeCache,
		userLimiter:    userLimiter,
		globalLimiter:  globalLimiter,
		promptTemplate: promptTemplate,
		deadline:       deadline,
	}
}

func (usecase GetNarrativeUseCase) Execute(ctx context.Context, userID int64) (string, error) {
	if !usecase.userLimiter.Allow(userID) {
		log.Printf("[narrative] user=%d rate_limit=exceeded", userID)
		return "", domain.ErrRateLimitExceeded
	}

	if cachedText, found := usecase.cache.Get(userID); found {
		log.Printf("[narrative] user=%d cache=hit", userID)
		return cachedText, nil
	}

	return usecase.runWithDeadline(userID)
}

func (usecase GetNarrativeUseCase) runWithDeadline(userID int64) (string, error) {
	outcomeChannel := make(chan narrativeOutcome, 1)
	startedAt := time.Now()

	go usecase.processInBackground(userID, outcomeChannel)

	select {
	case outcome := <-outcomeChannel:
		return usecase.handleOutcome(userID, outcome, startedAt), nil
	case <-time.After(usecase.deadline):
		log.Printf("[narrative] user=%d deadline=%dms action=processing_bg", userID, usecase.deadline.Milliseconds())
		return processingText, nil
	}
}

func (usecase GetNarrativeUseCase) handleOutcome(userID int64, outcome narrativeOutcome, startedAt time.Time) string {
	elapsedMs := time.Since(startedAt).Milliseconds()

	if outcome.err != nil {
		return usecase.translateError(userID, outcome.err, elapsedMs)
	}
	if !validateNarrative(outcome.text, outcome.tickers) {
		log.Printf("[narrative] user=%d total_ms=%d valid=false action=fallback", userID, elapsedMs)
		return fallbackText
	}
	log.Printf("[narrative] user=%d total_ms=%d cache=miss valid=true", userID, elapsedMs)
	return outcome.text
}

func (usecase GetNarrativeUseCase) translateError(userID int64, err error, elapsedMs int64) string {
	if errors.Is(err, errEmptyPortfolio) {
		return emptyPortfolioText
	}
	if errors.Is(err, errGlobalLimit) {
		log.Printf("[narrative] user=%d global_limit=exceeded action=fallback", userID)
		return fallbackText
	}
	log.Printf("[narrative] user=%d total_ms=%d err=%v action=fallback", userID, elapsedMs, err)
	return fallbackText
}

func (usecase GetNarrativeUseCase) processInBackground(userID int64, outcomeChannel chan<- narrativeOutcome) {
	snapshot, tickers, err := usecase.aggregate(userID)
	if err != nil {
		sendOutcome(outcomeChannel, narrativeOutcome{err: err})
		return
	}
	if len(tickers) == 0 {
		sendOutcome(outcomeChannel, narrativeOutcome{err: errEmptyPortfolio})
		return
	}
	if !usecase.globalLimiter.Allow() {
		sendOutcome(outcomeChannel, narrativeOutcome{err: errGlobalLimit})
		return
	}

	prompt, err := buildPrompt(usecase.promptTemplate, snapshot)
	if err != nil {
		sendOutcome(outcomeChannel, narrativeOutcome{err: err})
		return
	}

	generatedText, err := usecase.callLLM(prompt)
	if err != nil {
		log.Printf("[narrative] user=%d bg_llm_err=%v", userID, err)
	} else if !validateNarrative(generatedText, tickers) {
		log.Printf("[narrative] user=%d bg_valid=false", userID)
	} else {
		usecase.cache.Set(userID, generatedText, narrativeCacheTTL)
		log.Printf("[narrative] user=%d bg_cache=set", userID)
	}
	sendOutcome(outcomeChannel, narrativeOutcome{text: generatedText, tickers: tickers, err: err})
}

func (usecase GetNarrativeUseCase) callLLM(prompt string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), narrativeLLMTimeout)
	defer cancel()
	return usecase.llm.GenerateNarrative(ctx, prompt)
}

func sendOutcome(outcomeChannel chan<- narrativeOutcome, outcome narrativeOutcome) {
	select {
	case outcomeChannel <- outcome:
	default:
	}
}
