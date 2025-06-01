package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

const localsUserKey = "userID"

var (
	ErrUnauthorized = keyauth.ErrMissingOrMalformedAPIKey
)

func JWTMiddleware(service Service) fiber.Handler {
	return keyauth.New(keyauth.Config{
		KeyLookup: "cookie:token",
		Validator: func(c *fiber.Ctx, token string) (bool, error) {
			userID, err := service.ValidateToken(token)
			if err != nil {
				return false, keyauth.ErrMissingOrMalformedAPIKey
			}

			c.Locals(localsUserKey, userID)

			return true, nil
		},
	})
}

func GetUserID(c *fiber.Ctx) (string, bool) {
	userID, ok := c.Locals(localsUserKey).(string)
	return userID, ok
}
