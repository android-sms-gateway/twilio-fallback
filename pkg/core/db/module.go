package db

import (
	"database/sql"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"db",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("db")
	}),
	fx.Provide(func(cfg Config) (*sql.DB, error) {
		db, err := New(cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to open database: %w", err)
		}

		db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
		db.SetMaxOpenConns(cfg.MaxOpenConns)
		db.SetMaxIdleConns(cfg.MaxIdleConns)

		return db, nil
	}),
)
