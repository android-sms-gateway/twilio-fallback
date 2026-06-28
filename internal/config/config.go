package config

import (
	"fmt"
	"os"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/go-core-fx/config"
)

type httpConfig struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type twilioConfig struct {
	AccountSID  string `koanf:"account_sid"  required:"true"`
	AuthToken   string `koanf:"auth_token"   required:"true"`
	CallbackURL string `koanf:"callback_url"`
}

type smsGateConfig struct {
	BaseURL  string `koanf:"base_url"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`

	Timeout time.Duration `koanf:"timeout"`
}

type Config struct {
	HTTP    httpConfig    `koanf:"http"`
	Twilio  twilioConfig  `koanf:"twilio"`
	SMSGate smsGateConfig `koanf:"smsgate"`
}

func Default() Config {
	return Config{
		HTTP: httpConfig{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},
		Twilio: twilioConfig{
			AccountSID:  "",
			AuthToken:   "",
			CallbackURL: "",
		},
		SMSGate: smsGateConfig{
			BaseURL:  smsgateway.BaseURL,
			Timeout:  time.Second,
			Username: "",
			Password: "",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
