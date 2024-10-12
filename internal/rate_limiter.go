package internal

// RateLimiter is an interface that requires an Allow method.
// The Allow method will return a boolean value indicating
// whether the request is allowed or not based on the rate limiting strategy.
type RateLimiter interface {
	Allow() bool
}
