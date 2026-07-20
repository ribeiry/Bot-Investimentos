package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"portifolio-api/internal/domain"
	"strconv"
	"time"
)

type alphaVantageProvider struct {
	client *http.Client
	apiKey string
}

type alphaVantageResponse struct {
	GlobalQuote struct {
		Price         string `json:"05. price"`
		PreviousClose string `json:"08. previous close"`
	} `json:"Global Quote"`
}

func NewAlphaVantageProvider(token string) *alphaVantageProvider {
	return &alphaVantageProvider{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: token,
	}
}

func (a alphaVantageProvider) GetByTickers(assets []domain.Asset) ([]domain.Quote, error) {
	var quotes []domain.Quote

	for _, asset := range assets {
		if asset.Market == "B3" {
			continue // B3 é responsabilidade da Brapi
		}
		quote, err := a.GetByTicker(asset.Ticker, asset.Market)
		if err != nil {
			continue
		}
		quotes = append(quotes, *quote)
	}
	return quotes, nil
}

func (a alphaVantageProvider) GetByTicker(ticker string, market string) (*domain.Quote, error) {
	var alphaRespose alphaVantageResponse

	// URL base sem sufixo (NYSE)
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s&apikey=%s", ticker, a.apiKey)

	// se B3, adiciona .SA
	if market == "B3" {
		url = fmt.Sprintf("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=%s.SA&apikey=%s", ticker, a.apiKey)
	}
	resp, error := a.client.Get(url)

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&alphaRespose)

	price, _ := strconv.ParseFloat(alphaRespose.GlobalQuote.Price, 64)
	previousClose, _ := strconv.ParseFloat(alphaRespose.GlobalQuote.PreviousClose, 64)
	return &domain.Quote{
		Ticker:        ticker,
		CurrentValue:  price,
		PreviousValue: previousClose,
		QuotedAt:      time.Now(),
	}, nil

}
