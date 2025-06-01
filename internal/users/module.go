package users

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"users",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("users")
	}),
	fx.Provide(
		NewRepository,
		fx.Private,
	),
	fx.Provide(
		NewService,
	),
)
