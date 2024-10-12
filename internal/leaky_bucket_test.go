package internal

import (
	"testing"
	"time"
)

func TestLeakyBucketAllow(t *testing.T) {
	limiter := NewLeakyBucket(3, time.Millisecond*500)

	// Allow the first 3 requests
	for i := 0; i < 3; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d should have been allowed", i+1)
		}
	}

	// The 4th request should be denied
	if limiter.Allow() {
		t.Error("Request 4 should have been denied")
	}

	// Wait for 1 second, which should leak 2 requests
	time.Sleep(time.Second)

	// Now the next 2 requests should be allowed
	for i := 0; i < 2; i++ {
		if !limiter.Allow() {
			t.Errorf("Request after leak %d should have been allowed", i+1)
		}
	}

	// The next request should be denied
	if limiter.Allow() {
		t.Error("Request after leak should have been denied")
	}
}
