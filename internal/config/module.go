package config

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/encryption"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	fx.Provide(New),
	fx.Provide(func(c Config) http.Config {
		return http.Config{
			Address:     c.Http.Address,
			ProxyHeader: c.Http.ProxyHeader,
			Proxies:     c.Http.Proxies,
		}
	}),
	fx.Provide(func(c Config) encryption.Config {
		return encryption.Config{
			Key: c.Encryption.Key,
		}
	}),
	fx.Provide(func(c Config) auth.Config {
		return auth.Config{
			Secret: c.Auth.Secret,
			Expiry: c.Auth.Expiry,
		}
	}),
)
