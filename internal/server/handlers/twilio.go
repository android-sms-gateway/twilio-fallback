package handlers

import (
	"fmt"

	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/twilio"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TwilioHandler struct {
	handler.Base

	proxyService  proxy.Service
	twilioService twilio.Service

	logger *zap.Logger
}

func NewTwilioHandler(
	proxy proxy.Service,
	twilio twilio.Service,
	validator *validator.Validate,
	logger *zap.Logger,
) (handler.Handler, error) {
	if proxy == nil {
		return nil, ErrProxyServiceNil
	}

	if twilio == nil {
		return nil, ErrTwilioServiceNil
	}

	if validator == nil {
		return nil, ErrValidatorNil
	}

	if logger == nil {
		return nil, ErrLoggerNil
	}

	return &TwilioHandler{
		Base: handler.Base{
			Validator: validator,
		},
		proxyService:  proxy,
		twilioService: twilio,

		logger: logger,
	}, nil
}

// Register sets up the Twilio callback route.
func (h *TwilioHandler) Register(router fiber.Router) {
	router.Post("/twilio", h.callback)
}

//	@Summary		Handle Twilio message status callback
//	@Description	Processes Twilio message status updates with signature validation
//	@Tags			Twilio
//	@Accept			x-www-form-urlencoded
//	@Param			MessageSid			formData	string	true	"Twilio message SID"
//	@Param			MessageStatus		formData	string	true	"Message status (e.g. delivered, failed)"
//	@Param			X-Twilio-Signature	header		string	true	"Twilio request signature"
//	@Success		200					{string}	string	"OK"
//	@Failure		400					{string}	string	"Bad request"
//	@Failure		500					{string}	string	"Internal server error"
//	@Router			/api/twilio [post]
//
// Handle Twilio status callbacks.
func (h *TwilioHandler) callback(c *fiber.Ctx) error {
	scheme := "https"
	if c.Protocol() == "http" {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, c.Hostname(), c.OriginalURL())

	params := map[string]string{}
	for key, value := range c.Request().PostArgs().All() {
		params[string(key)] = string(value)
	}

	signature := c.Get("X-Twilio-Signature")
	if signature == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing X-Twilio-Signature header")
	}

	if err := h.twilioService.ValidateSignature(url, params, signature); err != nil {
		h.logger.Error("Invalid signature", zap.Error(err))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid signature")
	}

	// Parse form data and validate required fields
	messageSid := params["MessageSid"]
	messageStatus := params["MessageStatus"]

	// Validate required fields
	if messageSid == "" {
		return fiber.NewError(fiber.StatusBadRequest, "MessageSid is required")
	}
	if messageStatus == "" {
		return fiber.NewError(fiber.StatusBadRequest, "MessageStatus is required")
	}

	// Delegate processing to proxy service

	if err := h.proxyService.ProcessCallback(c.Context(), messageSid, messageStatus); err != nil {
		h.logger.Error("Failed to process callback", zap.String("message_sid", messageSid), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process callback: %s", err))
	}

	h.logger.Info("Callback processed", zap.String("message_sid", messageSid))

	return c.SendStatus(fiber.StatusOK)
}
