package server

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"server",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("server")
	}),
	fx.Provide(http.NewJSONErrorHandler),
	fx.Provide(func(log *zap.Logger) http.Options {
		return *(&http.Options{}).WithErrorHandler(http.NewJSONErrorHandler(log))
	}),
	fx.Provide(handlers.NewHealthHandler, fx.Private),
	fx.Provide(handlers.NewTwilioHandler, fx.Private),
	fx.Invoke(func(app *fiber.App, hh *handlers.HealthHandler, th *handlers.TwilioHandler) {
		hh.Register(app.Group("/health"))

		api := app.Group("/api")
		th.Register(api.Group("/twilio"))
	}),
)
