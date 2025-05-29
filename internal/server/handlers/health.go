package handlers

import (
	"fmt"

	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/health"
	"github.com/android-sms-gateway/twilio-fallback/internal/version"
	"github.com/capcom6/go-helpers/maps"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthHandler struct {
	handler.Base

	service *health.Service
}

//	@Summary		Health check
//	@Description	Checks if service is healthy
//	@Tags			System
//	@Produce		json
//	@Success		200	{object}	smsgateway.HealthResponse	"Health check result"
//	@Failure		503	{object}	smsgateway.HealthResponse	"Service is unhealthy"
//	@Router			/health [get]
//
// Health check
func (h *HealthHandler) get(c *fiber.Ctx) error {
	check, err := h.service.HealthCheck(c.Context())
	if err != nil {
		return fmt.Errorf("health check: %w", err)
	}

	res := smsgateway.HealthResponse{
		Status:    smsgateway.HealthStatus(check.Status),
		Version:   version.AppVersion,
		ReleaseID: version.AppReleaseID(),
		Checks: maps.MapValues(
			check.Checks,
			func(c health.CheckDetail) smsgateway.HealthCheck {
				return smsgateway.HealthCheck{
					Description:   c.Description,
					ObservedUnit:  c.ObservedUnit,
					ObservedValue: c.ObservedValue,
					Status:        smsgateway.HealthStatus(c.Status),
				}
			},
		),
	}

	if check.Status == health.StatusFail {
		return c.Status(fiber.StatusServiceUnavailable).JSON(res)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (c *HealthHandler) Register(router fiber.Router) {
	router.Get("", c.get)
}

func NewHealthHandler(service *health.Service, v *validator.Validate, l *zap.Logger) *HealthHandler {
	return &HealthHandler{
		Base: handler.Base{
			Validator: v,
			Logger:    l,
		},
		service: service,
	}
}
