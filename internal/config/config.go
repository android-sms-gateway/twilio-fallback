package config

import (
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
	DSN string `envconfig:"DATABASE__DSN"`
}

type EncryptionConfig struct {
	Key string `envconfig:"ENCRYPTION__KEY"`
}

type AuthConfig struct {
	Secret string        `envconfig:"AUTH__SECRET"`
	Expiry time.Duration `envconfig:"AUTH__EXPIRY"`
}

type Config struct {
	Http       HttpConfig
	Redis      RedisConfig
	Database   DatabaseConfig
	Encryption EncryptionConfig
	Auth       AuthConfig
}

var instance = Config{
	Http: HttpConfig{
		Address: "127.0.0.1:3000",
	},
	Redis: RedisConfig{
		URL: "redis://localhost:6379",
	},
	Database: DatabaseConfig{
		DSN: "mysql://root@tcp(localhost:3306)/twilio-fallback?charset=utf8mb4&parseTime=True&loc=Local",
	},
	Encryption: EncryptionConfig{},
	Auth: AuthConfig{
		Expiry: 24 * time.Hour,
	},
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
