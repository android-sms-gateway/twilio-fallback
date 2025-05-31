package orm

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"orm",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("orm")
	}),
	fx.Provide(New),
)
