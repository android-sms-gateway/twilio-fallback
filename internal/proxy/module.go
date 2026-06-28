package proxy

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"proxy",
		logger.WithNamedLogger("proxy"),

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
}
