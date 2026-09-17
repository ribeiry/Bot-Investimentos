package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow(userID int64) bool
}

type dailyCounter struct {
	requestCount int
	resetAt      time.Time
}

type dailyLimiter struct {
	mutex        sync.Mutex
	dailyLimit   int
	countersByUser map[int64]dailyCounter
}

func NewDailyLimiter(dailyLimit int) *dailyLimiter {
	return &dailyLimiter{
		dailyLimit:     dailyLimit,
		countersByUser: map[int64]dailyCounter{},
	}
}

func (limiter *dailyLimiter) Allow(userID int64) bool {
	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := time.Now().UTC()
	counter := limiter.countersByUser[userID]

	if now.After(counter.resetAt) {
		counter = dailyCounter{requestCount: 0, resetAt: nextMidnightUTC(now)}
	}
	if counter.requestCount >= limiter.dailyLimit {
		limiter.countersByUser[userID] = counter
		return false
	}
	counter.requestCount++
	limiter.countersByUser[userID] = counter
	return true
}

func nextMidnightUTC(now time.Time) time.Time {
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(24 * time.Hour)
}
