package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Test TokenBucket Middleware
func TestTokenBucketMiddleware(t *testing.T) {
	middleware, err := RateLimitingMiddleware("token_bucket", 2, time.Second)
	if err != nil {
		t.Fatal("Error creating middleware:", err)
	}

	// Create a sample handler that returns a simple response
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	}))

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Send the first request (should be allowed)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	// Send the second request (should be allowed)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	// Send the third request (should be rate-limited)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status TooManyRequests, got %d", rr.Code)
	}

	// Wait for the bucket to refill and try again
	time.Sleep(time.Second)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK after refill, got %d", rr.Code)
	}
}

// Test LeakyBucket Middleware
func TestLeakyBucketMiddleware(t *testing.T) {
	middleware, err := RateLimitingMiddleware("leaky_bucket", 2, time.Millisecond*500)
	if err != nil {
		t.Fatal("Error creating middleware:", err)
	}

	// Create a sample handler that returns a simple response
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	}))

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Send the first request (should be allowed)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	// Send the second request (should be allowed)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	// Send the third request (should be rate-limited)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status TooManyRequests, got %d", rr.Code)
	}

	// Wait for the bucket to "leak" and try again
	time.Sleep(time.Second)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK after leak, got %d", rr.Code)
	}
}
