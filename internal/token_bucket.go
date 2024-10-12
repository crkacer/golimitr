package internal

import (
	"sync"
	"time"
)

type TokenBucket struct {
	limit     int
	remaining int
	interval  time.Duration
	lastCheck time.Time
	mu        sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastCheck)
	tokensToAdd := int(elapsed / tb.interval)

	if tokensToAdd > 0 {
		tb.remaining += tokensToAdd
		if tb.remaining > tb.limit {
			tb.remaining = tb.limit
		}
		tb.lastCheck = now
	}

	if tb.remaining > 0 {
		tb.remaining--
		return true
	}
	return false
}

func NewTokenBucket(limit int, interval time.Duration) RateLimiter {
	return &TokenBucket{
		limit:     limit,
		remaining: limit,
		interval:  interval,
		lastCheck: time.Now(),
	}
}
