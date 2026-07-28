package repository

import (
	"database/sql"
	"portifolio-api/internal/domain"
	"time"
)

type priceHistoryRepository struct {
	db *sql.DB
}

func NewPriceHistoryRepository(db *sql.DB) *priceHistoryRepository {
	return &priceHistoryRepository{db: db}
}

func (r priceHistoryRepository) Save(asset domain.PriceHistory) error {

	query := `INSERT INTO price_history (ticker, price, captured_at) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, asset.Ticker, asset.Price, asset.CapturedAt)
	return err
}

func (r priceHistoryRepository) GetLastPrice(ticker string) (float64, error) {
	var price float64
	err := r.db.QueryRow(
		"SELECT price FROM price_history WHERE ticker = ? ORDER BY captured_at DESC LIMIT 1",
		ticker,
	).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (r priceHistoryRepository) GetPriceAtDate(ticker string, date time.Time) (float64, error) {
	var price float64
	err := r.db.QueryRow(
		`SELECT price FROM price_history
		 WHERE ticker = ? AND captured_at <= ?
		 ORDER BY captured_at DESC LIMIT 1`,
		ticker, date,
	).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
