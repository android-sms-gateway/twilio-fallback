package smsgate

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"smsgate",
		logger.WithNamedLogger("smsgate"),

		fx.Provide(NewService),
	)
}
