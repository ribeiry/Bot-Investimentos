package ratelimit

import (
	"sync"
	"time"
)

type GlobalRateLimiter interface {
	Allow() bool
}

type globalLimiter struct {
	mutex           sync.Mutex
	limitPerMinute  int
	limitPerDay     int
	minuteCount     int
	dayCount        int
	minuteResetAt   time.Time
	dayResetAt      time.Time
}

func NewGlobalLimiter(limitPerMinute, limitPerDay int) *globalLimiter {
	now := time.Now().UTC()
	return &globalLimiter{
		limitPerMinute: limitPerMinute,
		limitPerDay:    limitPerDay,
		minuteResetAt:  now.Add(time.Minute),
		dayResetAt:     nextMidnightUTC(now),
	}
}

func (limiter *globalLimiter) Allow() bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := time.Now().UTC()
	if now.After(limiter.minuteResetAt) {
		limiter.minuteCount = 0
		limiter.minuteResetAt = now.Add(time.Minute)
	}
	if now.After(limiter.dayResetAt) {
		limiter.dayCount = 0
		limiter.dayResetAt = nextMidnightUTC(now)
	}
	if limiter.minuteCount >= limiter.limitPerMinute || limiter.dayCount >= limiter.limitPerDay {
		return false
	}
	limiter.minuteCount++
	limiter.dayCount++
	return true
}
