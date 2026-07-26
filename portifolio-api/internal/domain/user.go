package domain

import "time"

type User struct {
	ID         int64
	TelegramID string
	Name       string
	APIKey     string
	CreatedAt  time.Time
}

type UserRepository interface {
	FindByAPIKey(apiKey string) (*User, error)
	Create(user User) (*User, error)
}
