package config

import (
	"github.com/android-sms-gateway/twilio-fallback/internal/smsgate"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(New, fx.Private),
		fx.Provide(
			func(c Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     c.HTTP.Address,
					ProxyHeader: c.HTTP.ProxyHeader,
					Proxies:     c.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
		),

		fx.Provide(func(c Config) twilio.Config {
			return twilio.Config{
				AccountSID:  c.Twilio.AccountSID,
				AuthToken:   c.Twilio.AuthToken,
				CallbackURL: c.Twilio.CallbackURL,
			}
		}),
		fx.Provide(func(c Config) smsgate.Config {
			return smsgate.Config{
				BaseURL:  c.SMSGate.BaseURL,
				Username: c.SMSGate.Username,
				Password: c.SMSGate.Password,

				Timeout: c.SMSGate.Timeout,
			}
		}),
	)
}
