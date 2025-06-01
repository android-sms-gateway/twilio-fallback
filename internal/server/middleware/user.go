package middleware

import (
	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const localsUserKey = "user"

func UserMiddleware(userService users.Service, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from locals (set by JWTMiddleware)
		userID, ok := auth.GetUserID(c)
		if !ok {
			logger.Warn("user ID not found in locals")
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}

		// Get user from database
		user, err := userService.GetUser(userID)
		if err != nil {
			if users.IsUserNotFound(err) {
				logger.Warn("user not found in database", zap.String("user_id", userID))
				return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
			}
			logger.Error("error getting user from database", zap.Error(err))
			return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
		}

		// Set user in locals
		c.Locals(localsUserKey, user)

		return c.Next()
	}
}

func GetUser(c *fiber.Ctx) (*users.User, bool) {
	user, ok := c.Locals(localsUserKey).(*users.User)
	return user, ok
}

func WithUser(fn func(*fiber.Ctx, *users.User) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := GetUser(c)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
		}
		return fn(c, user)
	}
}
