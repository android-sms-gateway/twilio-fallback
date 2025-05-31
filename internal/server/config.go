package server

import "time"

type RateLimitConfig struct {
	Requests int           `envconfig:"SERVER__RATE_LIMIT__REQUESTS"`
	Period   time.Duration `envconfig:"SERVER__RATE_LIMIT__PERIOD"`
}

type Config struct {
	RateLimit RateLimitConfig
}
