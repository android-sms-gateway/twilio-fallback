package internal

import (
	"github.com/android-sms-gateway/twilio-fallback/internal/config"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server"
	"github.com/android-sms-gateway/twilio-fallback/internal/smsgate"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/validatorfx"
	"go.uber.org/fx"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// badgerfx.Module(),
		// bunfx.Module(),
		// cachefx.Module(),
		fiberfx.Module(),
		// gocqlfx.Module(),
		// gocqlxfx.Module(),
		// sqlfx.Module(),
		// goosefx.Module(),
		// gormfx.Module(),
		healthfx.Module(),
		// openrouterfx.Module(),
		// redisfx.Module(),
		// sqlxfx.Module(),
		// telegofx.Module(true),
		validatorfx.Module(),
		// watermillfx.Module(),
		//
		// APP MODULES
		config.Module(),
		twilio.Module(),
		smsgate.Module(),
		server.Module(),
		//
		// BUSINESS MODULES
		fx.Supply(version),
		proxy.Module(),
	).
		Run()
}
