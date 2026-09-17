package mocks

import (
	mock "github.com/stretchr/testify/mock"
)

type RateLimiter struct {
	mock.Mock
}

func (m *RateLimiter) Allow(userID int64) bool {
	ret := m.Called(userID)
	return ret.Bool(0)
}

type GlobalRateLimiter struct {
	mock.Mock
}

func (m *GlobalRateLimiter) Allow() bool {
	ret := m.Called()
	return ret.Bool(0)
}
