package orm

import (
	"database/sql"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

func New(config Config, db *sql.DB, logger *zap.Logger) (*gorm.DB, error) {
	log := zapgorm2.New(logger)
	log.LogLevel = glogger.Info
	if config.Debug {
		log.LogLevel += 1
	}
	// log.SetAsDefault()

	return gorm.Open(
		mysql.New(mysql.Config{
			Conn: db,
		}),
		&gorm.Config{
			Logger: log,
		},
	)
}
