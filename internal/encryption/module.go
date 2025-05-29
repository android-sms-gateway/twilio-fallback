package encryption

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"encryption",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("encryption")
	}),
	fx.Provide(
		func(c Config) (Encryptor, error) {
			return NewEncryptor([]byte(c.Key))
		},
	),
)
