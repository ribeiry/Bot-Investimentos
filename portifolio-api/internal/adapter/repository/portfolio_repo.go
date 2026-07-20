package repository

import (
	"database/sql"
	"portifolio-api/internal/domain"
)

type portfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *portfolioRepository {
	return &portfolioRepository{db: db}
}

func (r portfolioRepository) Upsert(asset domain.Asset) error {

	query := `INSERT INTO portfolio (ticker,market,quantity,average_price) VALUES (?,?,?,?)
	ON CONFLICT (TICKER) DO UPDATE SET QUANTITY = EXCLUDED.QUANTITY, AVERAGE_PRICE = EXCLUDED.AVERAGE_PRICE`

	_, err := r.db.Exec(query, asset.Ticker, asset.Market, asset.Quantity, asset.AveragePrice)
	return err
}

func (r portfolioRepository) ReturnAllPortfolio() ([]domain.Asset, error) {
	var assets []domain.Asset
	query := "SELECT * FROM PORTFOLIO;"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var asset domain.Asset
		var id int
		rows.Scan(&id, &asset.Ticker, &asset.Market, &asset.Quantity, &asset.AveragePrice, &asset.CreatedAt)
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r portfolioRepository) ReturnAssetPortfolio(ticker string) (*domain.Asset, error) {

	var asset domain.Asset

	query := "SELECT * FROM PORTFOLIO WHERE TICKER = ?;"
	row := r.db.QueryRow(query, ticker)
	error := row.Scan(&asset.Ticker, &asset.Market, &asset.Quantity, &asset.AveragePrice)
	if error != nil {
		return nil, error

	}

	return &asset, nil

}

func (r portfolioRepository) DeleteByTicker(ticker string) error {

	query := "DELETE FROM PORTFOLIO WHERE TICKER = ?;"

	_, error := r.db.Exec(query, ticker)

	return error

}
