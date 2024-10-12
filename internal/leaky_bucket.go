// The Leaky Bucket algorithm allows requests to flow at a consistent rate by "leaking" them over time.
// Excess requests are discarded once the bucket capacity is exceeded.
package internal

import (
	"sync"
	"time"
)

// LeakyBucket implements a thread-safe leaky bucket rate limiter.
// The bucket has a defined capacity and rate of leakage. The water level (requests) is reduced
// over time, simulating the steady leak of water from a bucket.
type LeakyBucket struct {
	capacity int           // The maximum number of requests the bucket can hold.
	rate     time.Duration // The rate at which the bucket leaks (one unit per rate interval).
	water    int           // Current water level (number of requests in the bucket).
	lastLeak time.Time     // The time when the last leak occurred.
	mu       sync.Mutex    // Mutex to ensure thread-safe access to the bucket.
}

// Allow checks whether a new request is allowed to enter the bucket.
// The method calculates how much water has leaked since the last request
// and adjusts the water level accordingly. If the bucket is not full, the request is allowed
// and the water level is increased. If the bucket is full, the request is rejected.
//
// Returns true if the request is allowed, false if it is denied.
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.lastLeak)  // Calculate the time elapsed since the last leak
	leaked := int(elapsed / lb.rate) // Determine how much water has leaked in that time

	// If enough time has passed, reduce the water level by the leaked amount
	if leaked > 0 {
		lb.water -= leaked
		// Ensure the water level doesn't drop below zero
		if lb.water < 0 {
			lb.water = 0
		}
		lb.lastLeak = now
	}

	// If the current water level is below the bucket's capacity, allow the request
	if lb.water < lb.capacity {
		lb.water++ // Add one unit of water for the new request
		return true
	}
	// Otherwise, reject the request as the bucket is full
	return false
}

// NewLeakyBucket creates and returns a new instance of a LeakyBucket rate limiter.
// It takes two parameters: the bucket's capacity and the rate at which it leaks.
// The capacity defines the maximum number of requests that can be processed before
// the bucket becomes full, and the rate determines how quickly the bucket "leaks" and allows new requests.
//
// Parameters:
// - capacity: The maximum number of requests that can be processed in a burst before limiting.
// - rate: The time duration between each leak event.
//
// Returns a RateLimiter interface which can be used to control the rate of requests.
func NewLeakyBucket(capacity int, rate time.Duration) RateLimiter {
	return &LeakyBucket{
		capacity: capacity,   // Set the maximum capacity of the bucket
		rate:     rate,       // Set the rate at which the bucket leaks
		lastLeak: time.Now(), // Initialize the lastLeak time to the current time
	}
}
