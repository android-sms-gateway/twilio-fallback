package config

import (
	"runtime"
	"time"

	"github.com/android-sms-gateway/core/config"
)

type HttpConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type RedisConfig struct {
	URL string `envconfig:"REDIS__URL"`
}

type DatabaseConfig struct {
	URL             string        `envconfig:"DATABASE__URL"`
	ConnMaxIdleTime time.Duration `envconfig:"DATABASE__CONN_MAX_IDLE_TIME"`
	ConnMaxLifetime time.Duration `envconfig:"DATABASE__CONN_MAX_LIFETIME"`
	MaxOpenConns    int           `envconfig:"DATABASE__MAX_OPEN_CONNS"`
	MaxIdleConns    int           `envconfig:"DATABASE__MAX_IDLE_CONNS"`
}

type EncryptionConfig struct {
	Key string `envconfig:"ENCRYPTION__KEY"`
}

type AuthConfig struct {
	Secret string        `envconfig:"AUTH__SECRET"`
	Expiry time.Duration `envconfig:"AUTH__EXPIRY"`
}

type RateLimitConfig struct {
	Requests int           `envconfig:"SERVER__RATE_LIMIT__REQUESTS"`
	Period   time.Duration `envconfig:"SERVER__RATE_LIMIT__PERIOD"`
}

type ServerConfig struct {
	RateLimit RateLimitConfig
}

type ProxyConfig struct {
	BaseURL string `envconfig:"PROXY__BASE_URL"`
}

type Config struct {
	Http       HttpConfig
	Redis      RedisConfig
	Database   DatabaseConfig
	Encryption EncryptionConfig
	Auth       AuthConfig
	Server     ServerConfig
	Proxy      ProxyConfig
	Debug      bool `envconfig:"DEBUG"`
}

var instance = Config{
	Http: HttpConfig{
		Address: "127.0.0.1:3000",
	},
	Redis: RedisConfig{
		URL: "redis://localhost:6379",
	},
	Database: DatabaseConfig{
		URL:             "mysql://root@tcp(localhost:3306)/twilio-fallback?charset=utf8mb4&parseTime=True&loc=Local",
		ConnMaxIdleTime: 3 * time.Minute,
		ConnMaxLifetime: 30 * time.Minute,
		MaxOpenConns:    runtime.NumCPU() * 4,
		MaxIdleConns:    runtime.NumCPU() * 2,
	},
	Encryption: EncryptionConfig{},
	Auth: AuthConfig{
		Expiry: 24 * time.Hour,
	},
	Server: ServerConfig{
		RateLimit: RateLimitConfig{
			Requests: 100,
			Period:   time.Minute,
		},
	},
	Proxy: ProxyConfig{
		BaseURL: "http://localhost:3000",
	},
	Debug: false,
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
