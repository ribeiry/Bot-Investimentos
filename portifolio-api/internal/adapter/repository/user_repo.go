package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"portifolio-api/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r userRepository) FindByAPIKey(apiKey string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, telegram_id, name, api_key, created_at FROM users WHERE api_key = $1;"
	row := r.db.QueryRow(query, apiKey)
	if err := row.Scan(&user.ID, &user.TelegramID, &user.Name, &user.APIKey, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r userRepository) FindByTelegramID(telegramID string) (*domain.User, error) {
	var user domain.User
	query := "SELECT id, telegram_id, name, api_key, created_at FROM users WHERE telegram_id = $1;"
	row := r.db.QueryRow(query, telegramID)
	if err := row.Scan(&user.ID, &user.TelegramID, &user.Name, &user.APIKey, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r userRepository) Create(user domain.User) (*domain.User, error) {
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, err
	}
	user.APIKey = apiKey

	query := "INSERT INTO users (telegram_id, name, api_key) VALUES ($1, $2, $3) RETURNING id;"
	var id int64
	if err := r.db.QueryRow(query, user.TelegramID, user.Name, user.APIKey).Scan(&id); err != nil {
		return nil, err
	}
	user.ID = id
	return &user, nil
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
