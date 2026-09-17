package domain

import (
	"context"
	"errors"
)

var ErrRateLimitExceeded = errors.New("limite diário atingido")

type LLMProvider interface {
	GenerateNarrative(ctx context.Context, prompt string) (string, error)
}
