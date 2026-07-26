package db

import "database/sql"

func RunMigrations(db *sql.DB) error {
	if err := createPortfolioTable(db); err != nil {
		return err
	}
	if err := createPriceHistoryTable(db); err != nil {
		return err
	}
	if err := createUsersTable(db); err != nil {
		return err
	}
	if err := migratePortfolioAddUserID(db); err != nil {
		return err
	}
	return migrateUsersAddAPIKey(db)
}

func createPortfolioTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS portfolio (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL DEFAULT 0,
		ticker VARCHAR(20) NOT NULL,
		market VARCHAR(10) NOT NULL,
		quantity DECIMAL(10,4) NOT NULL,
		average_price DECIMAL(10,2) NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(ticker, user_id)
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
		api_key VARCHAR(64) UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func migratePortfolioAddUserID(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(portfolio)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasUserID := false
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var dfltValue interface{}
		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
		if name == "user_id" {
			hasUserID = true
			break
		}
	}
	if hasUserID {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmts := []string{
		`CREATE TABLE portfolio_new (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL DEFAULT 0,
			ticker VARCHAR(20) NOT NULL,
			market VARCHAR(10) NOT NULL,
			quantity DECIMAL(10,4) NOT NULL,
			average_price DECIMAL(10,2) NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(ticker, user_id)
		);`,
		`INSERT INTO portfolio_new (id, user_id, ticker, market, quantity, average_price, created_at)
		 SELECT id, 0, ticker, market, quantity, average_price, created_at FROM portfolio;`,
		`DROP TABLE portfolio;`,
		`ALTER TABLE portfolio_new RENAME TO portfolio;`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func migrateUsersAddAPIKey(db *sql.DB) error {
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasAPIKey := false
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var dfltValue interface{}
		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
		if name == "api_key" {
			hasAPIKey = true
			break
		}
	}
	if hasAPIKey {
		return nil
	}

	_, err = db.Exec(`ALTER TABLE users ADD COLUMN api_key VARCHAR(64) UNIQUE;`)
	return err
}
