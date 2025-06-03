package handlers

import (
	"fmt"

	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TwilioHandler struct {
	handler.Base

	proxyService  proxy.Service
	twilioService twilio.Service
}

func NewTwilioHandler(proxy proxy.Service, twilio twilio.Service, validator *validator.Validate, logger *zap.Logger) *TwilioHandler {
	if proxy == nil {
		panic("proxy service is nil")
	}

	if twilio == nil {
		panic("twilio service is nil")
	}

	if validator == nil {
		panic("validator is nil")
	}

	if logger == nil {
		panic("logger is nil")
	}

	return &TwilioHandler{
		Base: handler.Base{
			Validator: validator,
			Logger:    logger,
		},
		proxyService:  proxy,
		twilioService: twilio,
	}
}

func (h *TwilioHandler) callback(c *fiber.Ctx) error {
	// Validate Twilio signature
	url := fmt.Sprintf("https://%s%s", c.Hostname(), c.OriginalURL())
	params := map[string]string{}
	if err := c.BodyParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Failed to parse request body: %s", err))
	}

	signature := c.Get("X-Twilio-Signature")
	if signature == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing X-Twilio-Signature header")
	}

	if err := h.twilioService.ValidateSignature(url, params, signature); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Twilio signature validation failed: %s", err))
	}

	// Parse form data
	messageSid := c.FormValue("MessageSid")
	messageStatus := c.FormValue("MessageStatus")

	// Delegate processing to proxy service

	if err := h.proxyService.ProcessCallback(c.Context(), messageSid, messageStatus); err != nil {
		h.Logger.Error("Failed to process callback", zap.String("message_sid", messageSid), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process callback: %s", err))
	}

	h.Logger.Info("Callback processed", zap.String("message_sid", messageSid))

	return c.SendStatus(fiber.StatusOK)
}

func (h *TwilioHandler) Register(router fiber.Router) {
	router.Post("", h.callback)
}
