package mocks

import (
	"time"

	mock "github.com/stretchr/testify/mock"
)

type NarrativeCache struct {
	mock.Mock
}

func (m *NarrativeCache) Get(userID int64) (string, bool) {
	ret := m.Called(userID)
	return ret.String(0), ret.Bool(1)
}

func (m *NarrativeCache) Set(userID int64, text string, ttl time.Duration) {
	m.Called(userID, text, ttl)
}
