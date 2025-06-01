package server

import (
	"github.com/android-sms-gateway/core/http"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/middleware"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/redis/go-redis/v9"
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
		return *(&http.Options{}).
			WithErrorHandler(http.NewJSONErrorHandler(log)).
			WithViews(html.New("internal/server/views", ".html"))
	}),
	fx.Provide(handlers.NewHealthHandler, fx.Private),
	fx.Provide(handlers.NewAuthHandler, fx.Private),
	fx.Provide(handlers.NewDashboardHandler, fx.Private),
	fx.Invoke(func(
		app *fiber.App,
		rdb *redis.Client,
		hh *handlers.HealthHandler,
		ah *handlers.AuthHandler,
		dh *handlers.DashboardHandler,
		authService auth.Service,
		userService users.Service,
		proxy proxy.Service,
		logger *zap.Logger,
		config Config,
	) {
		rateLimitMiddleware := middleware.NewRateLimit(rdb, config.RateLimit.Requests, config.RateLimit.Period)

		// Apply rate limiter to sensitive routes
		authGroup := app.Group("/auth", rateLimitMiddleware)
		ah.Register(authGroup)

		// Register other routes
		hh.Register(app.Group("/health"))

		// Apply JWT middleware
		jwtMiddleware := auth.JWTMiddleware(authService)

		// Apply user middleware
		userMiddleware := middleware.UserMiddleware(userService, logger)

		// Protect dashboard routes
		protected := app.Group("/dashboard", jwtMiddleware, userMiddleware)
		dh.Register(protected)
	}),
)
