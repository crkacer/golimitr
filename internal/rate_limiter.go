package internal

type RateLimiter interface {
	Allow() bool
}
