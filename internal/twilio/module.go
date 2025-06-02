package twilio

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"twilio",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("twilio")
	}),
	fx.Provide(NewService),
)
