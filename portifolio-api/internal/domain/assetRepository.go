package domain

type AssetRepository interface {
	Upsert(asset Asset) error
	ReturnAllPortfolio() ([]Asset, error)
	ReturnAssetPortfolio(ticker string) (*Asset, error)
	DeleteByTicker(ticker string) error
}
