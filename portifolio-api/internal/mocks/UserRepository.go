package mocks

import (
	domain "portifolio-api/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (_m *UserRepository) FindByAPIKey(apiKey string) (*domain.User, error) {
	ret := _m.Called(apiKey)

	if fn, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return fn(apiKey)
	}

	var user *domain.User
	if ret.Get(0) != nil {
		user = ret.Get(0).(*domain.User)
	}
	return user, ret.Error(1)
}

func (_m *UserRepository) FindByTelegramID(telegramID string) (*domain.User, error) {
	ret := _m.Called(telegramID)

	if fn, ok := ret.Get(0).(func(string) (*domain.User, error)); ok {
		return fn(telegramID)
	}

	var user *domain.User
	if ret.Get(0) != nil {
		user = ret.Get(0).(*domain.User)
	}
	return user, ret.Error(1)
}

func (_m *UserRepository) Create(user domain.User) (*domain.User, error) {
	ret := _m.Called(user)

	if fn, ok := ret.Get(0).(func(domain.User) (*domain.User, error)); ok {
		return fn(user)
	}

	var created *domain.User
	if ret.Get(0) != nil {
		created = ret.Get(0).(*domain.User)
	}
	return created, ret.Error(1)
}

func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	m := &UserRepository{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}
