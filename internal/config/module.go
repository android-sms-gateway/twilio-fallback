package config

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/twilio-fallback/internal/smsgate"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
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
	fx.Provide(func(c Config) twilio.Config {
		return twilio.Config{
			AccountSID: c.Twilio.AccountSID,
			AuthToken:  c.Twilio.AuthToken,
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
