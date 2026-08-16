package db

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"createUsersTable", createUsersTable},
		{"createPortfolioTable", createPortfolioTable},
		{"createPriceHistoryTable", createPriceHistoryTable},
		{"createAlertsTable", createAlertsTable},
		{"addUsersTelegramIDUniqueIndex", addUsersTelegramIDUniqueIndex},
	}

	for _, m := range migrations {
		log.Printf("[migrations] rodando %s", m.name)
		if err := m.fn(db); err != nil {
			log.Printf("[migrations] falha em %s: %v", m.name, err)
			return err
		}
	}
	return nil
}

func createUsersTable(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id          SERIAL PRIMARY KEY,
		telegram_id VARCHAR(50) NOT NULL,
		name        VARCHAR(100),
		api_key     VARCHAR(64),
		created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_api_key ON users(api_key);`)
	return err
}

func createPortfolioTable(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS portfolio (
		id            SERIAL PRIMARY KEY,
		user_id       INTEGER        NOT NULL DEFAULT 0,
		ticker        VARCHAR(20)    NOT NULL,
		market        VARCHAR(10)    NOT NULL,
		quantity      DECIMAL(10,4)  NOT NULL,
		average_price DECIMAL(10,2)  NOT NULL,
		sector        VARCHAR(50),
		created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_portfolio_user_ticker ON portfolio(user_id, ticker);`)
	return err
}

func createPriceHistoryTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS price_history (
		id          SERIAL PRIMARY KEY,
		ticker      VARCHAR(20)   NOT NULL,
		price       DECIMAL(10,2) NOT NULL,
		captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func addUsersTelegramIDUniqueIndex(db *sql.DB) error {
	_, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_id);`)
	return err
}

func createAlertsTable(db *sql.DB) error {
	if _, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS alerts (
		id         SERIAL PRIMARY KEY,
		user_id    INTEGER       NOT NULL,
		ticker     VARCHAR(20)   NOT NULL,
		market     VARCHAR(10)   NOT NULL,
		stop_gain  DECIMAL(10,2),
		stop_loss  DECIMAL(10,2),
		active     BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_alerts_user_ticker ON alerts(user_id, ticker);`)
	return err
}
