package domain

import (
	"errors"
	"time"
)

var ErrTelegramIDAlreadyExists = errors.New("telegram_id já cadastrado")

type User struct {
	ID         int64
	TelegramID string
	Name       string
	APIKey     string
	CreatedAt  time.Time
}

type UserRepository interface {
	FindByAPIKey(apiKey string) (*User, error)
	FindByTelegramID(telegramID string) (*User, error)
	Create(user User) (*User, error)
}
