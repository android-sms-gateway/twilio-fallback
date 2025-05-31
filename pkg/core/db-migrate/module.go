package dbmigrate

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"db-migrate",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("db-migrate")
	}),
	fx.Provide(fx.Annotate(New, fx.As(new(Migrator)))),
)
