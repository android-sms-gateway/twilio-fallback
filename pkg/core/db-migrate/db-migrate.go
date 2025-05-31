package dbmigrate

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

type Migrator interface {
	Up(ctx context.Context) error
	Down(ctx context.Context) error
}

type Service struct {
	db *sql.DB

	migrationsFS   fs.FS
	migrationsPath string
}

func New(db *sql.DB, config Config) *Service {
	return &Service{
		db: db,

		migrationsFS:   config.MigrationsFS,
		migrationsPath: config.MigrationsPath,
	}
}

// Up implements Migrator.
func (s *Service) Up(ctx context.Context) error {
	goose.SetBaseFS(s.migrationsFS)

	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Up(s.db, s.migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Down implements Migrator.
func (s *Service) Down(ctx context.Context) error {
	goose.SetBaseFS(s.migrationsFS)

	if err := goose.SetDialect(string(goose.DialectMySQL)); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}
	if err := goose.Down(s.db, s.migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

var _ = (Migrator)(new(Service))
