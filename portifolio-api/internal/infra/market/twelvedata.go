package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portifolio-api/internal/domain"
	"strconv"
	"strings"
	"time"
)

type twelveDataProvider struct {
	client *http.Client
	apiKey string
}

func NewTwelveDataProvider(apiKey string) *twelveDataProvider {
	return &twelveDataProvider{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: apiKey,
	}
}

type twelveDataMultiResponse map[string]struct {
	Close string `json:"close"`
}

type twelveDataSingleResponse struct {
	Symbol string `json:"symbol"`
	Close  string `json:"close"`
}

func (t twelveDataProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var usaAssets []domain.Asset
	for _, asset := range assets {
		if asset.Market != "B3" {
			usaAssets = append(usaAssets, asset)
		}
	}

	var quotes []domain.Quote
	batchSize := 8

	for i := 0; i < len(usaAssets); i += batchSize {
		end := i + batchSize
		if end > len(usaAssets) {
			end = len(usaAssets)
		}
		batch := usaAssets[i:end]

		// monta symbols do lote
		var symbols []string
		for _, a := range batch {
			symbols = append(symbols, a.Ticker)
		}

		url := fmt.Sprintf("https://api.twelvedata.com/quote?symbol=%s&apikey=%s",
			strings.Join(symbols, ","), t.apiKey)

		resp, err := t.client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if len(batch) == 1 {
			var single twelveDataSingleResponse
			json.NewDecoder(resp.Body).Decode(&single)
			if single.Symbol != "" && single.Close != "" {
				price, _ := strconv.ParseFloat(single.Close, 64)
				quotes = append(quotes, domain.Quote{
					Ticker:       single.Symbol,
					CurrentValue: price,
					QuotedAt:     time.Now(),
				})
			}
		} else {
			var batchResp twelveDataMultiResponse
			json.NewDecoder(resp.Body).Decode(&batchResp)
			for ticker, data := range batchResp {
				price, _ := strconv.ParseFloat(data.Close, 64)
				quotes = append(quotes, domain.Quote{
					Ticker:       ticker,
					CurrentValue: price,
					QuotedAt:     time.Now(),
				})
			}
		}

		// delay entre lotes para respeitar rate limit
		if end < len(usaAssets) {
			time.Sleep(61 * time.Second)
		}
	}

	return quotes, nil
}
