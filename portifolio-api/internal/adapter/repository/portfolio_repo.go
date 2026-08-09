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

func (r portfolioRepository) Upsert(userID int64, asset domain.Asset) error {
	query := `INSERT INTO portfolio (user_id, ticker, market, quantity, average_price, sector) VALUES ($1,$2,$3,$4,$5,$6)
	ON CONFLICT (user_id, ticker) DO UPDATE SET
		quantity = EXCLUDED.quantity,
		average_price = EXCLUDED.average_price,
		sector = COALESCE(NULLIF(EXCLUDED.sector, ''), portfolio.sector)`
	_, err := r.db.Exec(query, userID, asset.Ticker, asset.Market, asset.Quantity, asset.AveragePrice, asset.Sector)
	return err
}

func (r portfolioRepository) ReturnAllPortfolio(userID int64) ([]domain.Asset, error) {
	var assets []domain.Asset
	query := "SELECT ticker, market, quantity, average_price, COALESCE(sector, ''), created_at FROM portfolio WHERE user_id = $1;"
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var asset domain.Asset
		rows.Scan(&asset.Ticker, &asset.Market, &asset.Quantity, &asset.AveragePrice, &asset.Sector, &asset.CreatedAt)
		assets = append(assets, asset)
	}
	return assets, nil
}

func (r portfolioRepository) ReturnAssetPortfolio(userID int64, ticker string) (*domain.Asset, error) {
	var asset domain.Asset
	query := "SELECT ticker, market, quantity, average_price, COALESCE(sector, '') FROM portfolio WHERE user_id = $1 AND ticker = $2;"
	row := r.db.QueryRow(query, userID, ticker)
	if err := row.Scan(&asset.Ticker, &asset.Market, &asset.Quantity, &asset.AveragePrice, &asset.Sector); err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r portfolioRepository) DeleteByTicker(userID int64, ticker string) error {
	query := "DELETE FROM portfolio WHERE user_id = $1 AND ticker = $2;"
	_, err := r.db.Exec(query, userID, ticker)
	return err
}

func (r portfolioRepository) UpdateSector(userID int64, ticker string, sector string) error {
	_, err := r.db.Exec(
		"UPDATE portfolio SET sector = $1 WHERE user_id = $2 AND ticker = $3;",
		sector, userID, ticker,
	)
	return err
}
