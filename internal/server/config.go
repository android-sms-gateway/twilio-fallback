package server

import "time"

type RateLimitConfig struct {
	Requests int
	Period   time.Duration
}

type Config struct {
	RateLimit RateLimitConfig
}
