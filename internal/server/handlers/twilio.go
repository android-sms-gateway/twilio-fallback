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

func NewTwilioHandler(proxy proxy.Service, twilio twilio.Service, validator *validator.Validate, logger *zap.Logger) (*TwilioHandler, error) {
	if proxy == nil {
		return nil, fmt.Errorf("proxy service is nil")
	}

	if twilio == nil {
		return nil, fmt.Errorf("twilio service is nil")
	}

	if validator == nil {
		return nil, fmt.Errorf("validator is nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	return &TwilioHandler{
		Base: handler.Base{
			Validator: validator,
			Logger:    logger,
		},
		proxyService:  proxy,
		twilioService: twilio,
	}, nil
}

//	@Summary		Handle Twilio message status callback
//	@Description	Processes Twilio message status updates with signature validation
//	@Tags			Twilio
//	@Accept			x-www-form-urlencoded
//	@Param			MessageSid			formData	string	true	"Twilio message SID"
//	@Param			MessageStatus		formData	string	true	"Message status (e.g. delivered, failed)"
//	@Param			X-Twilio-Signature	header		string	true	"Twilio request signature"
//	@Success		200 {string}	string	"OK"
//	@Failure		400	{string}	string	"Bad request"
//	@Failure		500	{string}	string	"Internal server error"
//	@Router			/api/twilio [post]
//
// Handle Twilio status callbacks
func (h *TwilioHandler) callback(c *fiber.Ctx) error {
	scheme := "https"
	if c.Protocol() == "http" {
		scheme = "http"
	}
	url := fmt.Sprintf("%s://%s%s", scheme, c.Hostname(), c.OriginalURL())

	params := map[string]string{}
	c.Request().PostArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})

	signature := c.Get("X-Twilio-Signature")
	if signature == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing X-Twilio-Signature header")
	}

	if err := h.twilioService.ValidateSignature(url, params, signature); err != nil {
		h.Logger.Error("Invalid signature", zap.Error(err))
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
		h.Logger.Error("Failed to process callback", zap.String("message_sid", messageSid), zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to process callback: %s", err))
	}

	h.Logger.Info("Callback processed", zap.String("message_sid", messageSid))

	return c.SendStatus(fiber.StatusOK)
}

// Register sets up the Twilio callback route
func (h *TwilioHandler) Register(router fiber.Router) {
	router.Post("", h.callback)
}
