package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/android-sms-gateway/core/handler"
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
)

type AuthHandler struct {
	handler.Base

	authService auth.Service
	userService users.Service
}

func NewAuthHandler(logger *zap.Logger, validator *validator.Validate, authService auth.Service, userService users.Service) *AuthHandler {
	return &AuthHandler{
		Base: handler.Base{
			Validator: validator,
			Logger:    logger,
		},

		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := h.BodyParserValidator(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.userService.RegisterUser(req.Login, req.Password, req.TwilioAccountSID, req.TwilioAuthToken)
	if err != nil {
		if err == users.ErrUserAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "user already exists"})
		}
		h.Logger.Error("Failed to register user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to register user"})
	}

	// Set session cookie/JWT
	token, err := h.authService.GenerateToken(user.Login)
	if err != nil {
		h.Logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	// Set cookie
	h.setAuthCookie(c, token)

	// Redirect to dashboard
	return c.Redirect("/dashboard")
}

func (h *AuthHandler) login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := h.BodyParserValidator(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	user, err := h.userService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		if err == users.ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		h.Logger.Error("Failed to authenticate user", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to authenticate user"})
	}

	// Set session cookie/JWT
	token, err := h.authService.GenerateToken(user.Login)
	if err != nil {
		h.Logger.Error("Failed to generate token", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate token"})
	}

	// Set cookie
	h.setAuthCookie(c, token)

	// Redirect to dashboard
	return c.Redirect("/dashboard")
}

func (h *AuthHandler) setAuthCookie(c *fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})
}

func (h *AuthHandler) Register(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
}
