package market

import (
	"errors"
	"portifolio-api/internal/domain"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeProvider struct {
	mu         sync.Mutex
	callCount  int32
	delay      time.Duration
	returnQtes []domain.Quote
	returnErr  error
}

func (f *fakeProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	atomic.AddInt32(&f.callCount, 1)
	if f.delay > 0 {
		time.Sleep(f.delay)
	}
	if f.returnErr != nil {
		return nil, f.returnErr
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]domain.Quote(nil), f.returnQtes...), nil
}

func (f *fakeProvider) calls() int32 {
	return atomic.LoadInt32(&f.callCount)
}

type fakePriceHistory struct{}

func (fakePriceHistory) Save(ph domain.PriceHistory) error                     { return nil }
func (fakePriceHistory) GetLastPrice(_ string) (float64, error)                { return 0, nil }
func (fakePriceHistory) GetPriceAtDate(_ string, _ time.Time) (float64, error) { return 0, nil }

func TestBuildAssetsKey_OrdemEstavel(t *testing.T) {
	a := []domain.Asset{{Ticker: "MSFT", Market: "NASDAQ"}, {Ticker: "AAPL", Market: "NASDAQ"}, {Ticker: "PETR4", Market: "B3"}}
	b := []domain.Asset{{Ticker: "PETR4", Market: "B3"}, {Ticker: "AAPL", Market: "NASDAQ"}, {Ticker: "MSFT", Market: "NASDAQ"}}

	assert.Equal(t, buildAssetsKey(a), buildAssetsKey(b))
	assert.Equal(t, "AAPL:NASDAQ,MSFT:NASDAQ,PETR4:B3", buildAssetsKey(a))
}

func TestBuildAssetsKey_MarketFazParteDaChave(t *testing.T) {
	a := []domain.Asset{{Ticker: "VALE", Market: "B3"}}
	b := []domain.Asset{{Ticker: "VALE", Market: "NYSE"}}

	assert.NotEqual(t, buildAssetsKey(a), buildAssetsKey(b))
}

func TestGetByTickers_SingleFlight_DedupChamadasConcorrentes(t *testing.T) {
	brapi := &fakeProvider{
		delay:      100 * time.Millisecond,
		returnQtes: []domain.Quote{{Ticker: "PETR4", CurrentValue: 40}},
	}
	twelve := &fakeProvider{
		delay:      100 * time.Millisecond,
		returnQtes: []domain.Quote{{Ticker: "AAPL", CurrentValue: 300}},
	}
	provider := NewMarketProviderWithFallback(brapi, twelve, fakePriceHistory{})

	assets := []domain.Asset{
		{Ticker: "PETR4", Market: "B3"},
		{Ticker: "AAPL", Market: "NASDAQ"},
	}

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)
	results := make([][]domain.Quote, goroutines)
	errs := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx], errs[idx] = provider.GetByTickers(assets)
		}(i)
	}
	wg.Wait()

	// Todas devem ter sucesso e receber os mesmos quotes
	for i := 0; i < goroutines; i++ {
		assert.NoError(t, errs[i])
		assert.Len(t, results[i], 2)
	}

	// Apesar de 10 goroutines, cada provider externo foi chamado 1x só
	assert.Equal(t, int32(1), brapi.calls(), "brapi deveria ser chamado 1x (singleflight)")
	assert.Equal(t, int32(1), twelve.calls(), "twelvedata deveria ser chamado 1x (singleflight)")
}

func TestGetByTickers_SingleFlight_NaoCompartilhaChavesDiferentes(t *testing.T) {
	brapi := &fakeProvider{
		delay:      50 * time.Millisecond,
		returnQtes: []domain.Quote{{Ticker: "PETR4", CurrentValue: 40}},
	}
	twelve := &fakeProvider{
		delay:      50 * time.Millisecond,
		returnQtes: []domain.Quote{{Ticker: "AAPL", CurrentValue: 300}},
	}
	provider := NewMarketProviderWithFallback(brapi, twelve, fakePriceHistory{})

	assetsA := []domain.Asset{{Ticker: "PETR4", Market: "B3"}}
	assetsB := []domain.Asset{{Ticker: "AAPL", Market: "NASDAQ"}}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); provider.GetByTickers(assetsA) }()
	go func() { defer wg.Done(); provider.GetByTickers(assetsB) }()
	wg.Wait()

	assert.Equal(t, int32(1), brapi.calls())
	assert.Equal(t, int32(1), twelve.calls())
}

func TestGetByTickers_SingleFlight_LiberaAposConclusao(t *testing.T) {
	brapi := &fakeProvider{
		returnQtes: []domain.Quote{{Ticker: "PETR4", CurrentValue: 40}},
	}
	provider := NewMarketProviderWithFallback(brapi, &fakeProvider{}, fakePriceHistory{})
	assets := []domain.Asset{{Ticker: "PETR4", Market: "B3"}}

	_, err := provider.GetByTickers(assets)
	assert.NoError(t, err)
	_, err = provider.GetByTickers(assets)
	assert.NoError(t, err)

	// Chamadas sequenciais NÃO compartilham — 2 calls externas
	assert.Equal(t, int32(2), brapi.calls())
}

func TestGetByTickers_PropagaErro(t *testing.T) {
	brapi := &fakeProvider{returnErr: errors.New("brapi down")}
	twelve := &fakeProvider{returnQtes: []domain.Quote{{Ticker: "AAPL", CurrentValue: 300}}}
	provider := NewMarketProviderWithFallback(brapi, twelve, fakePriceHistory{})

	assets := []domain.Asset{
		{Ticker: "PETR4", Market: "B3"},
		{Ticker: "AAPL", Market: "NASDAQ"},
	}
	_, err := provider.GetByTickers(assets)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "brapi down")
}

func TestGetByTickers_AssetsVazio(t *testing.T) {
	provider := NewMarketProviderWithFallback(&fakeProvider{}, &fakeProvider{}, fakePriceHistory{})
	quotes, err := provider.GetByTickers(nil)
	assert.NoError(t, err)
	assert.Nil(t, quotes)
}
