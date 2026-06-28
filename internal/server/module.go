package server

import (
	"github.com/android-sms-gateway/twilio-fallback/internal/server/docs"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/health"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		logger.WithNamedLogger("server"),

		fx.Provide(func(log *zap.Logger) fiberfx.Options {
			opts := fiberfx.Options{}
			opts.WithErrorHandler(fiberfx.NewJSONErrorHandler(log))
			opts.WithMetrics()
			return opts
		}),
		fx.Supply(docs.SwaggerInfo),

		fx.Provide(
			health.NewHandler,
			openapi.NewHandler,
			fx.Private,
		),

		fx.Provide(
			fx.Annotate(handlers.NewTwilioHandler, fx.ResultTags(`group:"handlers"`)),
			fx.Private,
		),

		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, healthHandler *health.Handler, openapiHandler *openapi.Handler, app *fiber.App) {
					// Health endpoint
					healthHandler.Register(app)

					// Version 1 API group
					api := app.Group("/api")
					openapiHandler.Register(api.Group("/docs"))

					api.Use(validation.Middleware)

					for _, h := range handlers {
						h.Register(api)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
		// fx.Invoke(func(app *fiber.App, hh *handlers.HealthHandler, th *handlers.TwilioHandler) {
		// 	hh.Register(app.Group("/health"))

		// 	api := app.Group("/api")
		// 	th.Register(api.Group("/twilio"))
		// }),
	)
}
