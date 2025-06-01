package internal

import (
	"context"
	"io/fs"

	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/core/logger"
	"github.com/android-sms-gateway/core/redis"
	"github.com/android-sms-gateway/core/validator"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/config"
	"github.com/android-sms-gateway/twilio-fallback/internal/encryption"
	"github.com/android-sms-gateway/twilio-fallback/internal/health"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
	"github.com/android-sms-gateway/twilio-fallback/pkg/core/db"
	dbmigrate "github.com/android-sms-gateway/twilio-fallback/pkg/core/db-migrate"
	"github.com/android-sms-gateway/twilio-fallback/pkg/core/orm"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

var loggerModule = fx.Options(
	logger.Module,
	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
		logOption := fxevent.ZapLogger{Logger: logger}
		logOption.UseLogLevel(zapcore.DebugLevel)
		return &logOption
	}),
)

var appModule = fx.Options(
	health.Module,
	users.Module,
	auth.Module,
	encryption.Module,
	proxy.Module,
)

var httpModule = fx.Options(
	http.Module,
	validator.Module,
	server.Module,
)

var dbModule = fx.Options(
	db.Module,
	orm.Module,
)

func Run() {
	fx.New(
		loggerModule,

		redis.Module,

		httpModule,
		dbModule,

		config.Module,
		appModule,
	).
		Run()
}

func RunORMMigrations() {
	fx.New(
		loggerModule,

		dbModule,

		config.Module,
		appModule,

		fx.Invoke(func(db *gorm.DB, shutdown fx.Shutdowner) {
			if err := orm.RunMigrations(db); err != nil {
				panic(err)
			}

			shutdown.Shutdown()
		}),
	).Run()
}

func RunMigrations(migrationsFS fs.FS) {
	fx.New(
		loggerModule,

		dbModule,
		dbmigrate.Module,

		config.Module,
		fx.Provide(func() dbmigrate.Config {
			return dbmigrate.Config{
				MigrationsFS:   migrationsFS,
				MigrationsPath: "migrations",
			}
		}),
		fx.Invoke(func(migrator dbmigrate.Migrator, shutdown fx.Shutdowner) {
			migrator.Up(context.Background())

			shutdown.Shutdown()
		}),
	).Run()
}
