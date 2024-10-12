// The middleware limits the rate of incoming requests and rejects those
// exceeding the allowed rate.
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/crkacer/golimitr/internal"
)

// CreateRateLimiter initializes a rate limiter based on the provided strategy.
// The function supports two rate limiting strategies: "token_bucket" and "leaky_bucket".
// It takes three parameters: the rate limiting strategy as a string, the limit on requests,
// and the interval or rate at which tokens are refilled or requests are leaked.
//
// Parameters:
// - strat: The strategy to use for rate limiting ("token_bucket" or "leaky_bucket").
// - limit: The maximum number of requests allowed within the defined time frame.
// - interval: The duration between refills of tokens or the leak rate for requests.
//
// Returns:
// - internal.RateLimiter: The appropriate rate limiter based on the selected strategy.
// - error: An error is returned if an invalid strategy is provided.
func CreateRateLimiter(strat string, limit int, interval time.Duration) (internal.RateLimiter, error) {
	switch strat {
	case "token_bucket":
		return internal.NewTokenBucket(limit, interval), nil
	case "leaky_bucket":
		return internal.NewLeakyBucket(limit, interval), nil
	default:
		return nil, fmt.Errorf("unknown strategy: %s", strat)
	}
}

// RateLimitingMiddleware creates an HTTP middleware that applies rate limiting based on the selected strategy.
// It leverages the CreateRateLimiter function to choose between the "token_bucket" and "leaky_bucket" algorithms,
// then wraps the given HTTP handler with rate limiting logic. If the rate limiter denies a request due to too many requests,
// it responds with an HTTP 429 (Too Many Requests) status code.
//
// Parameters:
// - strategy: The rate limiting strategy to use ("token_bucket" or "leaky_bucket").
// - limit: The maximum number of allowed requests within the defined time frame.
// - intervalOrRate: The time duration between refills (for token bucket) or leak rate (for leaky bucket).
//
// Returns:
// - A middleware function that can be applied to an HTTP handler.
// - An error if an invalid strategy is provided or if the rate limiter fails to initialize.
func RateLimitingMiddleware(strategy string, limit int, intervalOrRate time.Duration) (func(http.Handler) http.Handler, error) {
	limiter, err := CreateRateLimiter(strategy, limit, intervalOrRate)
	if err != nil {
		return nil, err
	}

	// Middleware function that checks the rate limiter before allowing the request to proceed.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Deny the request if the rate limiter rejects it due to too many requests.
			if !limiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			// Allow the request to proceed to the next handler if the rate limiter permits it.
			next.ServeHTTP(w, r)
		})
	}, nil
}
