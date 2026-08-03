package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portifolio-api/internal/domain"
	"time"
)

type brapiProvider struct {
	client *http.Client
	token  string
}

type brapiResponse struct {
	Results []struct {
		Symbol                     string  `json:"symbol"`
		RegularMarketPrice         float64 `json:"regularMarketPrice"`
		RegularMarketPreviousClose float64 `json:"regularMarketPreviousClose"`
	} `json:"results"`
}

func NewBrapiProvider(token string) *brapiProvider {
	return &brapiProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  token,
	}
}

func (b brapiProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var quotes []domain.Quote

	for _, asset := range assets {
		url := fmt.Sprintf("https://brapi.dev/api/quote/%s?token=%s", asset.Ticker, b.token)

		resp, err := b.client.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var brapiResp brapiResponse
		json.NewDecoder(resp.Body).Decode(&brapiResp)

		if len(brapiResp.Results) == 0 {
			continue
		}

		result := brapiResp.Results[0]
		quotes = append(quotes, domain.Quote{
			Ticker:        result.Symbol,
			CurrentValue:  result.RegularMarketPrice,
			PreviousValue: result.RegularMarketPreviousClose,
			QuotedAt:      time.Now(),
		})
	}

	return quotes, nil
}
