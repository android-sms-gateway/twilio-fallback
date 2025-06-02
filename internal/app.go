package internal

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/core/logger"
	"github.com/android-sms-gateway/core/validator"
	"github.com/android-sms-gateway/twilio-fallback/internal/config"
	"github.com/android-sms-gateway/twilio-fallback/internal/health"
	"github.com/android-sms-gateway/twilio-fallback/internal/server"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	fx.New(
		logger.Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		http.Module,
		validator.Module,

		config.Module,
		twilio.Module,

		health.Module,

		server.Module,
	).
		Run()
}
