package user

import (
	"errors"
	"portifolio-api/internal/domain"
	"portifolio-api/internal/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)

	input := CreateUserInput{TelegramID: "123456", Name: "João"}
	created := &domain.User{ID: 1, TelegramID: "123456", Name: "João", APIKey: "abc123"}

	userRepo.On("Create", domain.User{TelegramID: "123456", Name: "João"}).Return(created, nil)

	usecase := NewCreateUserUseCase(userRepo)
	result, err := usecase.Execute(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "abc123", result.APIKey)
	userRepo.AssertExpectations(t)
}

func TestCreateUser_TelegramIDVazio(t *testing.T) {
	userRepo := new(mocks.UserRepository)

	usecase := NewCreateUserUseCase(userRepo)
	result, err := usecase.Execute(CreateUserInput{TelegramID: "", Name: "João"})

	assert.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertNotCalled(t, "Create")
}

func TestCreateUser_RepoError(t *testing.T) {
	userRepo := new(mocks.UserRepository)

	input := CreateUserInput{TelegramID: "123456", Name: "João"}
	userRepo.On("Create", domain.User{TelegramID: "123456", Name: "João"}).Return(nil, errors.New("db error"))

	usecase := NewCreateUserUseCase(userRepo)
	result, err := usecase.Execute(input)

	assert.Error(t, err)
	assert.Nil(t, result)
	userRepo.AssertExpectations(t)
}

func TestCreateUser_SemNome(t *testing.T) {
	userRepo := new(mocks.UserRepository)

	input := CreateUserInput{TelegramID: "999", Name: ""}
	created := &domain.User{ID: 2, TelegramID: "999", Name: "", APIKey: "key999"}

	userRepo.On("Create", domain.User{TelegramID: "999", Name: ""}).Return(created, nil)

	usecase := NewCreateUserUseCase(userRepo)
	result, err := usecase.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), result.ID)
	userRepo.AssertExpectations(t)
}
