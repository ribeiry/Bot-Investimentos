package repository

import (
	"database/sql"
	"portifolio-api/internal/domain"
)

type alertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *alertRepository {
	return &alertRepository{db: db}
}

func (r alertRepository) Upsert(alert domain.Alert) error {
	query := `
	INSERT INTO alerts (user_id, ticker, market, stop_gain, stop_loss, active)
	VALUES (?, ?, ?, ?, ?, 1)
	ON CONFLICT (user_id, ticker) DO UPDATE SET
		stop_gain = EXCLUDED.stop_gain,
		stop_loss = EXCLUDED.stop_loss,
		active = 1`
	_, err := r.db.Exec(query, alert.UserID, alert.Ticker, alert.Market, alert.StopGain, alert.StopLoss)
	return err
}

func (r alertRepository) GetAllByUserID(userID int64) ([]domain.Alert, error) {
	query := `SELECT id, user_id, ticker, market, stop_gain, stop_loss, active, created_at
	          FROM alerts WHERE user_id = ? AND active = 1;`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []domain.Alert
	for rows.Next() {
		var a domain.Alert
		err := rows.Scan(&a.ID, &a.UserID, &a.Ticker, &a.Market, &a.StopGain, &a.StopLoss, &a.Active, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

func (r alertRepository) DeleteByTicker(userID int64, ticker string) error {
	_, err := r.db.Exec("DELETE FROM alerts WHERE user_id = ? AND ticker = ?;", userID, ticker)
	return err
}
