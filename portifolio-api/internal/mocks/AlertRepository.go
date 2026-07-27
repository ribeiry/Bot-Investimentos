package mocks

import (
	domain "portifolio-api/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

type AlertRepository struct {
	mock.Mock
}

func (_m *AlertRepository) Upsert(alert domain.Alert) error {
	ret := _m.Called(alert)

	if fn, ok := ret.Get(0).(func(domain.Alert) error); ok {
		return fn(alert)
	}
	return ret.Error(0)
}

func (_m *AlertRepository) GetAllByUserID(userID int64) ([]domain.Alert, error) {
	ret := _m.Called(userID)

	if fn, ok := ret.Get(0).(func(int64) ([]domain.Alert, error)); ok {
		return fn(userID)
	}

	var alerts []domain.Alert
	if ret.Get(0) != nil {
		alerts = ret.Get(0).([]domain.Alert)
	}

	return alerts, ret.Error(1)
}

func (_m *AlertRepository) DeleteByTicker(userID int64, ticker string) error {
	ret := _m.Called(userID, ticker)

	if fn, ok := ret.Get(0).(func(int64, string) error); ok {
		return fn(userID, ticker)
	}
	return ret.Error(0)
}

func NewAlertRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlertRepository {
	m := &AlertRepository{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}
