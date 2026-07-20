package domain

type MarketProvider interface {
	GetByTickers(ticker []Asset) ([]Quote, error)
}
