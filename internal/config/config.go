package config

import (
	"github.com/android-sms-gateway/core/config"
)

type HttpConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type Config struct {
	Http HttpConfig
}

var instance = Config{
	Http: HttpConfig{
		Address: "127.0.0.1:3000",
	},
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
