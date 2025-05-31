package server

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"server",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("api")
	}),
	fx.Provide(http.NewJSONErrorHandler),
	fx.Provide(func(log *zap.Logger) http.Options {
		return *(&http.Options{}).WithErrorHandler(http.NewJSONErrorHandler(log))
	}),
	fx.Provide(handlers.NewHealthHandler, fx.Private),
	fx.Provide(handlers.NewAuthHandler, fx.Private),
	fx.Provide(handlers.NewDashboardHandler, fx.Private),
	fx.Provide(func(rdb *redis.Client, config Config) fiber.Handler {
		return middleware.NewRateLimit(rdb, config.RateLimit.Requests, config.RateLimit.Period)
	}),
	fx.Invoke(func(
		app *fiber.App,
		hh *handlers.HealthHandler,
		ah *handlers.AuthHandler,
		dh *handlers.DashboardHandler,
		authService *auth.AuthService,
		rateLimiter fiber.Handler,
	) {
		// Apply rate limiter to sensitive routes
		authGroup := app.Group("/auth", rateLimiter)
		ah.Register(authGroup)

		// Register other routes
		hh.Register(app.Group("/health"))
		protected := app.Group("/dashboard", auth.JWTMiddleware(authService))
		dh.Register(protected)
	}),
)
