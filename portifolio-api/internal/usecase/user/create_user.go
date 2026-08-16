package user

import (
	"errors"
	"portifolio-api/internal/domain"
)

type CreateUserUseCase struct {
	userRepo domain.UserRepository
}

func NewCreateUserUseCase(userRepo domain.UserRepository) CreateUserUseCase {
	return CreateUserUseCase{userRepo: userRepo}
}

type CreateUserInput struct {
	TelegramID string `json:"telegram_id"`
	Name       string `json:"name"`
}

func (u CreateUserUseCase) Execute(input CreateUserInput) (*domain.User, error) {
	if input.TelegramID == "" {
		return nil, errors.New("telegram_id is required")
	}
	existing, err := u.userRepo.FindByTelegramID(input.TelegramID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domain.ErrTelegramIDAlreadyExists
	}
	return u.userRepo.Create(domain.User{
		TelegramID: input.TelegramID,
		Name:       input.Name,
	})
}
