package db

import "database/sql"

func RunMigrations(db *sql.DB) error {
	if err := createPortfolioTable(db); err != nil {
		return err
	}
	if err := createPriceHistoryTable(db); err != nil {
		return err
	}
	return createUsersTable(db)
}

func createPortfolioTable(db *sql.DB) error {

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS portfolio (
		id INTEGER PRIMARY KEY,
		ticker VARCHAR(20) NOT NULL UNIQUE,
		market VARCHAR(10) NOT NULL,
		quantity DECIMAL(10,4) NOT NULL,
		average_price DECIMAL(10,2) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func createPriceHistoryTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS price_history (
		id INTEGER PRIMARY KEY,
		ticker VARCHAR(20) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		captured_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}
func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		telegram_id VARCHAR(50) NOT NULL,
		name VARCHAR(100),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}
