package db

import "time"

type Config struct {
	URL string

	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
}
