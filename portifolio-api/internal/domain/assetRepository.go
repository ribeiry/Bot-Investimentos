package domain

type AssetRepository interface {
	Upsert(userID int64, asset Asset) error
	ReturnAllPortfolio(userID int64) ([]Asset, error)
	ReturnAssetPortfolio(userID int64, ticker string) (*Asset, error)
	DeleteByTicker(userID int64, ticker string) error
	UpdateSector(userID int64, ticker string, sector string) error
}
