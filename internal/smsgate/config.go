package smsgate

import "time"

type Config struct {
	BaseURL  string
	Username string
	Password string

	Timeout time.Duration
}
