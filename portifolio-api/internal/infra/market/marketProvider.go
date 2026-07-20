package market

import (
	"encoding/json"
	"log"
	"net/http"
	"portifolio-api/internal/domain"
	"strings"
	"time"
)

type yahooProvider struct {
	client *http.Client
}

type yahooBatchResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol                     string  `json:"symbol"`
			RegularMarketPrice         float64 `json:"regularMarketPrice"`
			RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
		} `json:"result"`
	} `json:"quoteResponse"`
}

func NewYahooProvider() *yahooProvider {
	return &yahooProvider{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (y yahooProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var symbols []string
	var quotes []domain.Quote
	var yahooBatchResp yahooBatchResponse

	for _, asset := range assets {
		if asset.Market == "B3" {
			symbols = append(symbols, asset.Ticker+".SA")
		} else {
			symbols = append(symbols, asset.Ticker)
		}
	}

	log.Println(symbols)

	url := "https://query1.finance.yahoo.com/v7/finance/quote?symbols=" + strings.Join(symbols, ",")

	resp, error := y.client.Get(url)

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&yahooBatchResp)
	log.Printf("[Yahoo Batch] symbols: %s", strings.Join(symbols, ","))
	log.Printf("[Yahoo Batch] results: %d", len(yahooBatchResp.QuoteResponse.Result))

	for _, result := range yahooBatchResp.QuoteResponse.Result {
		// remove o .SA para manter o ticker limpo
		ticker := strings.TrimSuffix(result.Symbol, ".SA")
		quotes = append(quotes, domain.Quote{
			Ticker:        ticker,
			CurrentValue:  result.RegularMarketPrice,
			PreviousValue: result.RegularMarketPreviousClose,
			QuotedAt:      time.Now(),
		})
	}

	return quotes, nil
}
