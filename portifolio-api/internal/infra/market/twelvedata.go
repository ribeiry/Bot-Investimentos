package market

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

		batchQuotes, err := t.fetchBatch(url, batch)
		if err != nil {
			return nil, err
		}
		quotes = append(quotes, batchQuotes...)

		// delay entre lotes para respeitar rate limit
		if end < len(usaAssets) {
			time.Sleep(61 * time.Second)
		}
	}

	return quotes, nil
}

func (t twelveDataProvider) fetchBatch(url string, batch []domain.Asset) ([]domain.Quote, error) {
	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("twelvedata http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("twelvedata read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[twelvedata] status=%d body=%s", resp.StatusCode, truncate(string(body), 200))
		return nil, fmt.Errorf("twelvedata http %d", resp.StatusCode)
	}

	if len(batch) == 1 {
		return parseSingleQuote(body)
	}
	return parseBatchQuotes(body)
}

func parseSingleQuote(body []byte) ([]domain.Quote, error) {
	var single twelveDataSingleResponse
	if err := json.Unmarshal(body, &single); err != nil {
		return nil, fmt.Errorf("twelvedata decode single: %w (body=%s)", err, truncate(string(body), 200))
	}
	if single.Symbol == "" || single.Close == "" {
		return nil, fmt.Errorf("twelvedata empty single response: %s", truncate(string(body), 200))
	}
	price, err := strconv.ParseFloat(single.Close, 64)
	if err != nil {
		return nil, fmt.Errorf("twelvedata parse price %q: %w", single.Close, err)
	}
	return []domain.Quote{{
		Ticker:       single.Symbol,
		CurrentValue: price,
		QuotedAt:     time.Now(),
	}}, nil
}

func parseBatchQuotes(body []byte) ([]domain.Quote, error) {
	var batchResp twelveDataMultiResponse
	if err := json.Unmarshal(body, &batchResp); err != nil {
		return nil, fmt.Errorf("twelvedata decode batch: %w (body=%s)", err, truncate(string(body), 200))
	}
	if len(batchResp) == 0 {
		return nil, fmt.Errorf("twelvedata empty batch response: %s", truncate(string(body), 200))
	}
	quotes := make([]domain.Quote, 0, len(batchResp))
	for ticker, data := range batchResp {
		price, err := strconv.ParseFloat(data.Close, 64)
		if err != nil {
			return nil, fmt.Errorf("twelvedata parse price ticker=%s value=%q: %w", ticker, data.Close, err)
		}
		quotes = append(quotes, domain.Quote{
			Ticker:       ticker,
			CurrentValue: price,
			QuotedAt:     time.Now(),
		})
	}
	return quotes, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
