// The Token Bucket algorithm allows a certain number of requests (tokens) to be processed
// within a given time frame, with tokens being refilled periodically at a set interval.
package internal

import (
	"sync"
	"time"
)

// TokenBucket implements a thread-safe token bucket rate limiter.
// The bucket has a limit on how many tokens it can hold, and tokens are consumed as requests are processed.
// Tokens are replenished over time at a fixed interval.
type TokenBucket struct {
	limit     int           // The maximum number of tokens the bucket can hold.
	remaining int           // The current number of remaining tokens.
	interval  time.Duration // The interval at which tokens are added back to the bucket.
	lastCheck time.Time     // The time of the last token refill check.
	mu        sync.Mutex    // Mutex to ensure thread-safe access to the bucket.
}

// Allow checks whether a new request is allowed based on the current number of tokens.
// If enough time has passed since the last check, tokens are refilled into the bucket.
// If there are tokens available, the request is allowed and a token is consumed.
// Otherwise, the request is denied.
//
// Returns true if the request is allowed (token is available), false if denied.
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastCheck)          // Calculate time elapsed since the last check
	tokensToAdd := int(elapsed / tb.interval) // Determine how many tokens to add based on the interval

	// Refill tokens based on the elapsed time
	if tokensToAdd > 0 {
		tb.remaining += tokensToAdd
		// Ensure the bucket does not exceed the maximum token limit
		if tb.remaining > tb.limit {
			tb.remaining = tb.limit
		}
		tb.lastCheck = now
	}

	// If there are tokens remaining, allow the request
	if tb.remaining > 0 {
		tb.remaining-- // Consume one token
		return true
	}
	// If no tokens are available, deny the request
	return false
}

// NewTokenBucket creates and returns a new instance of a TokenBucket rate limiter.
// It takes two parameters: the token limit and the interval at which tokens are refilled.
// The limit defines the maximum number of tokens the bucket can hold,
// and the interval determines how frequently tokens are added.
//
// Parameters:
// - limit: The maximum number of tokens that can be stored in the bucket.
// - interval: The duration between refills of tokens.
//
// Returns a RateLimiter interface that can be used to control the rate of requests.
func NewTokenBucket(limit int, interval time.Duration) RateLimiter {
	return &TokenBucket{
		limit:     limit,      // Set the maximum number of tokens
		remaining: limit,      // Initialize the bucket with the full amount of tokens
		interval:  interval,   // Set the interval at which tokens are refilled
		lastCheck: time.Now(), // Initialize the lastCheck time to the current time
	}
}
