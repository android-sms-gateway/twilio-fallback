package handlers

import (
	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/middleware"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	handler.Base

	proxy proxy.Service
}

func NewDashboardHandler(proxy proxy.Service, logger *zap.Logger, validator *validator.Validate) *DashboardHandler {
	return &DashboardHandler{
		Base: handler.Base{
			Logger:    logger,
			Validator: validator,
		},

		proxy: proxy,
	}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx, user *users.User) error {
	// Get callback URL
	callbackURL := h.proxy.GetCallbackURL(user.CallbackUUID)

	// Render template with callback URL
	return c.Render("dashboard", fiber.Map{
		"Message":     "Welcome to the dashboard",
		"User":        user.Login,
		"CallbackURL": callbackURL,
	})
}

func (h *DashboardHandler) Register(router fiber.Router) {
	router.Get("/", middleware.WithUser(h.GetDashboard))
}
