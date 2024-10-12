package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/crkacer/golimitr/internal"
)

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

func RateLimitingMiddleware(strategy string, limit int, intervalOrRate time.Duration) (func(http.Handler) http.Handler, error) {
	limiter, err := CreateRateLimiter(strategy, limit, intervalOrRate)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}, nil
}
