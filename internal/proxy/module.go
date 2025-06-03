package proxy

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"proxy",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("proxy")
	}),
	fx.Provide(NewService),
	fx.Invoke(func(lifecycle fx.Lifecycle, service Service) {
		// Register shutdown hook
		lifecycle.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				service.Start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				service.Stop()
				return nil
			},
		})
	}),
)
