package handlers

import (
	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type DashboardHandler struct {
	handler.Base
}

func NewDashboardHandler(logger *zap.Logger) *DashboardHandler {
	return &DashboardHandler{
		Base: handler.Base{
			Logger: logger,
		},
	}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	userID, ok := auth.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	h.Logger.Info("User accessed dashboard", zap.String("userID", userID))
	return c.JSON(fiber.Map{"message": "Welcome to the dashboard", "user": userID})
}

func (h *DashboardHandler) Register(router fiber.Router) {
	router.Get("/", h.GetDashboard)
}
