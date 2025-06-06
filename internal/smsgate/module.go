package smsgate

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"smsgate",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("smsgate")
	}),
	fx.Provide(NewService),
)
