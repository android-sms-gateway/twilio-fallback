package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
)

type AuthHandler struct {
	handler.Base

	authService *auth.AuthService
}

func NewAuthHandler(logger *zap.Logger, validator *validator.Validate, authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		Base: handler.Base{
			Validator: validator,
			Logger:    logger,
		},

		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := h.BodyParserValidator(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Mock user validation - replace with actual user validation
	if req.Username != "admin" || req.Password != "password" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := h.authService.GenerateToken(req.Username)
	if err != nil {
		h.Logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Register(router fiber.Router) {
	router.Post("/login", h.Login)
}
