package cache

import (
	"sync"
	"time"
)

type NarrativeCache interface {
	Get(userID int64) (string, bool)
	Set(userID int64, text string, ttl time.Duration)
}

type cacheEntry struct {
	text      string
	expiresAt time.Time
}

type inMemoryNarrativeCache struct {
	mutex   sync.RWMutex
	entries map[int64]cacheEntry
}

func NewNarrativeCache() *inMemoryNarrativeCache {
	return &inMemoryNarrativeCache{entries: map[int64]cacheEntry{}}
}

func (cache *inMemoryNarrativeCache) Get(userID int64) (string, bool) {
	cache.mutex.RLock()
	entry, found := cache.entries[userID]
	cache.mutex.RUnlock()
	if !found {
		return "", false
	}
	if time.Now().After(entry.expiresAt) {
		cache.mutex.Lock()
		delete(cache.entries, userID)
		cache.mutex.Unlock()
		return "", false
	}
	return entry.text, true
}

func (cache *inMemoryNarrativeCache) Set(userID int64, text string, ttl time.Duration) {
	cache.mutex.Lock()
	cache.entries[userID] = cacheEntry{text: text, expiresAt: time.Now().Add(ttl)}
	cache.mutex.Unlock()
}
