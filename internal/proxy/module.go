package proxy

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"proxy",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("proxy")
	}),
	fx.Provide(NewService),
)
