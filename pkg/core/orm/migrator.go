package orm

import (
	"gorm.io/gorm"
)

type Migrator func(*gorm.DB) error

var migrations = []Migrator{}

func RegisterMigration(migrator Migrator) {
	migrations = append(migrations, migrator)
}

func RunMigrations(db *gorm.DB) error {
	for _, migrator := range migrations {
		if err := migrator(db); err != nil {
			return err
		}
	}

	return nil
}
