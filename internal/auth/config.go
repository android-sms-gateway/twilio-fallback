package auth

import "time"

type Config struct {
	Secret string
	Expiry time.Duration
}
