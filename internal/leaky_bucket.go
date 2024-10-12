package internal

import (
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity int
	rate     time.Duration
	water    int
	lastLeak time.Time
	mu       sync.Mutex
}

func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak)
	leaked := int(elapsed / lb.rate)

	if leaked > 0 {
		lb.water -= leaked
		if lb.water < 0 {
			lb.water = 0
		}
		lb.lastLeak = now
	}

	if lb.water < lb.capacity {
		lb.water++
		return true
	}
	return false
}

func NewLeakyBucket(capacity int, rate time.Duration) RateLimiter {
	return &LeakyBucket{
		capacity: capacity,
		rate:     rate,
		lastLeak: time.Now(),
	}
}
