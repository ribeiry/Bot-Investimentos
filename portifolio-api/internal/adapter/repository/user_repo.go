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
	query := "SELECT id, telegram_id, name, api_key, created_at FROM users WHERE api_key = ?;"
	row := r.db.QueryRow(query, apiKey)
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

	query := "INSERT INTO users (telegram_id, name, api_key) VALUES (?, ?, ?);"
	result, err := r.db.Exec(query, user.TelegramID, user.Name, user.APIKey)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
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
