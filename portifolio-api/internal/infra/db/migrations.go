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
	if err := createAlertsTable(db); err != nil {
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
		id           INTEGER PRIMARY KEY,
		ticker       VARCHAR(20)    NOT NULL UNIQUE,
		market       VARCHAR(10)    NOT NULL,
		quantity     DECIMAL(10,4)  NOT NULL,
		average_price DECIMAL(10,2) NOT NULL,
		created_at   DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func createPriceHistoryTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS price_history (
		id          INTEGER PRIMARY KEY,
		ticker      VARCHAR(20)   NOT NULL,
		price       DECIMAL(10,2) NOT NULL,
		captured_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id          INTEGER PRIMARY KEY,
		telegram_id VARCHAR(50) NOT NULL,
		name        VARCHAR(100),
		created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func createAlertsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS alerts (
		id         INTEGER PRIMARY KEY,
		user_id    INTEGER       NOT NULL,
		ticker     VARCHAR(20)   NOT NULL,
		market     VARCHAR(10)   NOT NULL,
		stop_gain  DECIMAL(10,2),
		stop_loss  DECIMAL(10,2),
		active     BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_alerts_user_ticker ON alerts(user_id, ticker);`)
	return err
}

// migratePortfolioAddUserID adds user_id to portfolio and replaces the
// UNIQUE(ticker) constraint with UNIQUE(user_id, ticker).
// SQLite does not support ALTER TABLE ADD CONSTRAINT, so we recreate the table.
func migratePortfolioAddUserID(db *sql.DB) error {
	if hasColumn(db, "portfolio", "user_id") {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmts := []string{
		`DROP TABLE IF EXISTS portfolio_new;`,
		`CREATE TABLE portfolio_new (
			id            INTEGER PRIMARY KEY,
			user_id       INTEGER       NOT NULL DEFAULT 0,
			ticker        VARCHAR(20)   NOT NULL,
			market        VARCHAR(10)   NOT NULL,
			quantity      DECIMAL(10,4) NOT NULL,
			average_price DECIMAL(10,2) NOT NULL,
			created_at    DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`INSERT INTO portfolio_new (id, user_id, ticker, market, quantity, average_price, created_at)
		 SELECT id, 0, ticker, market, quantity, average_price, created_at FROM portfolio;`,
		`DROP TABLE portfolio;`,
		`ALTER TABLE portfolio_new RENAME TO portfolio;`,
		`CREATE UNIQUE INDEX idx_portfolio_user_ticker ON portfolio(user_id, ticker);`,
	}

	for _, stmt := range stmts {
		if _, err := tx.Exec(stmt); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// migrateUsersAddAPIKey adds the api_key column to users.
// SQLite does not support UNIQUE in ALTER TABLE ADD COLUMN,
// so we add the column then create a separate unique index.
func migrateUsersAddAPIKey(db *sql.DB) error {
	if hasColumn(db, "users", "api_key") {
		return nil
	}

	if _, err := db.Exec(`ALTER TABLE users ADD COLUMN api_key VARCHAR(64);`); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_api_key ON users(api_key);`)
	return err
}

func hasColumn(db *sql.DB, table, column string) bool {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notNull, pk int
		var name, colType string
		var dfltValue interface{}
		rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
		if name == column {
			return true
		}
	}
	return false
}
