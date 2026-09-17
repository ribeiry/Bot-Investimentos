package mocks

import (
	"context"

	mock "github.com/stretchr/testify/mock"
)

type LLMProvider struct {
	mock.Mock
}

func (m *LLMProvider) GenerateNarrative(ctx context.Context, prompt string) (string, error) {
	ret := m.Called(ctx, prompt)

	if fn, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return fn(ctx, prompt)
	}
	return ret.String(0), ret.Error(1)
}
