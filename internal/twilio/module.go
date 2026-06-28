package twilio

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"twilio",
		logger.WithNamedLogger("twilio"),

		fx.Provide(NewService),
	)
}
