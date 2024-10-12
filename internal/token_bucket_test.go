package internal

import (
	"testing"
	"time"
)

func TestTokenBucketAllow(t *testing.T) {
	limiter := NewTokenBucket(5, time.Second)

	// Allow the first 5 requests
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// The 6th request should be denied
	if limiter.Allow() {
		t.Error("Request 6 should have been denied")
	}

	// Wait for 1 second for the bucket to refill
	time.Sleep(time.Second)

	// Now the next request should be allowed
	if !limiter.Allow() {
		t.Error("Request after refill should have been allowed")
	}
}
