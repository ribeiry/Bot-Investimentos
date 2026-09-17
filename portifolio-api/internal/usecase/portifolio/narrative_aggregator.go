package portifolio

import (
	"log"
	"portifolio-api/internal/domain"
	"sync"
)

type portfolioTotals struct {
	Invested         float64 `json:"invested_value"`
	Current          float64 `json:"current_value"`
	ProfitLoss       float64 `json:"profit_loss"`
	ReturnPercentage float64 `json:"return_percentage"`
}

type portfolioSnapshot struct {
	Totals      portfolioTotals `json:"totals"`
	Summary     any             `json:"summary"`
	Performance any             `json:"performance"`
	Benchmark   any             `json:"benchmark"`
	Allocation  any             `json:"allocation"`
}

func computeTotals(performance []domain.AssetPerformance) portfolioTotals {
	totals := portfolioTotals{}
	for _, asset := range performance {
		totals.Invested += asset.InvestedValue
		totals.Current += asset.CurrentValue
	}
	totals.ProfitLoss = totals.Current - totals.Invested
	if totals.Invested > 0 {
		totals.ReturnPercentage = totals.ProfitLoss / totals.Invested * 100
	}
	return totals
}

type aggregationState struct {
	mutex    sync.Mutex
	snapshot portfolioSnapshot
	tickers  []string
	firstErr error
}

func (state *aggregationState) recordError(err error) {
	state.mutex.Lock()
	defer state.mutex.Unlock()
	if state.firstErr == nil {
		state.firstErr = err
	}
}

func (state *aggregationState) setSummary(summary any) {
	state.mutex.Lock()
	state.snapshot.Summary = summary
	state.mutex.Unlock()
}

func (state *aggregationState) setBenchmark(benchmark any) {
	state.mutex.Lock()
	state.snapshot.Benchmark = benchmark
	state.mutex.Unlock()
}

func (state *aggregationState) setAllocation(allocation any) {
	state.mutex.Lock()
	state.snapshot.Allocation = allocation
	state.mutex.Unlock()
}

func (state *aggregationState) setPerformance(performance []domain.AssetPerformance, tickers []string) {
	state.mutex.Lock()
	state.snapshot.Performance = performance
	state.snapshot.Totals = computeTotals(performance)
	state.tickers = tickers
	state.mutex.Unlock()
}

func (usecase GetNarrativeUseCase) aggregate(userID int64) (portfolioSnapshot, []string, error) {
	state := &aggregationState{}
	var waitGroup sync.WaitGroup

	waitGroup.Add(4)
	go usecase.loadSummary(userID, state, &waitGroup)
	go usecase.loadPerformance(userID, state, &waitGroup)
	go usecase.loadBenchmark(userID, state, &waitGroup)
	go usecase.loadAllocation(userID, state, &waitGroup)
	waitGroup.Wait()

	if state.firstErr != nil {
		return state.snapshot, nil, state.firstErr
	}
	return state.snapshot, state.tickers, nil
}

func (usecase GetNarrativeUseCase) loadSummary(userID int64, state *aggregationState, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	summary, err := usecase.getSummary.Execute(userID, "cached", "market")
	if err != nil {
		state.recordError(err)
		return
	}
	state.setSummary(summary)
}

func (usecase GetNarrativeUseCase) loadPerformance(userID int64, state *aggregationState, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	performance, err := usecase.getPerformance.Execute(userID)
	if err != nil {
		state.recordError(err)
		return
	}
	tickers := make([]string, 0, len(performance))
	for _, asset := range performance {
		tickers = append(tickers, asset.Ticker)
	}
	state.setPerformance(performance, tickers)
}

func (usecase GetNarrativeUseCase) loadBenchmark(userID int64, state *aggregationState, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	benchmark, err := usecase.getBenchmark.Execute(userID, "weekly")
	if err != nil {
		log.Printf("[narrative] user=%d benchmark_err=%v (não-fatal)", userID, err)
		return
	}
	state.setBenchmark(benchmark)
}

func (usecase GetNarrativeUseCase) loadAllocation(userID int64, state *aggregationState, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	allocation, err := usecase.getAllocation.Execute(userID)
	if err != nil {
		state.recordError(err)
		return
	}
	state.setAllocation(allocation)
}
