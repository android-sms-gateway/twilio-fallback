package config

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/core/redis"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/encryption"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server"
	"github.com/android-sms-gateway/twilio-fallback/pkg/core/db"
	"github.com/android-sms-gateway/twilio-fallback/pkg/core/orm"
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
	fx.Provide(func(c Config) redis.Config {
		return redis.Config{
			URL: c.Redis.URL,
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
	fx.Provide(func(c Config) server.Config {
		return server.Config{
			RateLimit: server.RateLimitConfig{
				Requests: c.Server.RateLimit.Requests,
				Period:   c.Server.RateLimit.Period,
			},
		}
	}),
	fx.Provide(func(c Config) db.Config {
		return db.Config{
			URL:             c.Database.URL,
			ConnMaxIdleTime: c.Database.ConnMaxIdleTime,
			ConnMaxLifetime: c.Database.ConnMaxLifetime,
			MaxOpenConns:    c.Database.MaxOpenConns,
			MaxIdleConns:    c.Database.MaxIdleConns,
		}
	}),
	fx.Provide(func(c Config) orm.Config {
		return orm.Config{
			Debug: c.Debug,
		}
	}),
	fx.Provide(func(c Config) proxy.Config {
		return proxy.Config{
			BaseURL: c.Proxy.BaseURL,
		}
	}),
)
