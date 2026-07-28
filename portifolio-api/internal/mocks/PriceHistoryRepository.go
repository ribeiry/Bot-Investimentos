package mocks

import (
	domain "portifolio-api/internal/domain"
	"time"

	mock "github.com/stretchr/testify/mock"
)

type PriceHistoryRepository struct {
	mock.Mock
}

func (_m *PriceHistoryRepository) Save(history domain.PriceHistory) error {
	ret := _m.Called(history)

	if fn, ok := ret.Get(0).(func(domain.PriceHistory) error); ok {
		return fn(history)
	}
	return ret.Error(0)
}

func (_m *PriceHistoryRepository) GetLastPrice(ticker string) (float64, error) {
	ret := _m.Called(ticker)

	if fn, ok := ret.Get(0).(func(string) (float64, error)); ok {
		return fn(ticker)
	}
	return ret.Get(0).(float64), ret.Error(1)
}

func (_m *PriceHistoryRepository) GetPriceAtDate(ticker string, date time.Time) (float64, error) {
	ret := _m.Called(ticker, date)

	if fn, ok := ret.Get(0).(func(string, time.Time) (float64, error)); ok {
		return fn(ticker, date)
	}
	return ret.Get(0).(float64), ret.Error(1)
}

func NewPriceHistoryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *PriceHistoryRepository {
	m := &PriceHistoryRepository{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}
