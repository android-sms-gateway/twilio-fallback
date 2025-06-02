package config

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/core/config"
)

type HttpConfig struct {
	Address     string   `envconfig:"HTTP__ADDRESS"`
	ProxyHeader string   `envconfig:"HTTP__PROXY_HEADER"`
	Proxies     []string `envconfig:"HTTP__PROXIES"`
}

type TwilioConfig struct {
	AccountSID string `envconfig:"TWILIO__ACCOUNT_SID" required:"true"`
	AuthToken  string `envconfig:"TWILIO__AUTH_TOKEN" required:"true"`
}

type SMSGateConfig struct {
	BaseURL  string `envconfig:"SMSGATE__BASE_URL"`
	Username string `envconfig:"SMSGATE__USERNAME" required:"true"`
	Password string `envconfig:"SMSGATE__PASSWORD" required:"true"`
}

type Config struct {
	Http    HttpConfig
	Twilio  TwilioConfig
	SMSGate SMSGateConfig
}

var instance = Config{
	Http: HttpConfig{
		Address: "127.0.0.1:3000",
	},
	Twilio: TwilioConfig{},
	SMSGate: SMSGateConfig{
		BaseURL: smsgateway.BASE_URL,
	},
}

func New() (Config, error) {
	return instance, config.Load(&instance)
}
