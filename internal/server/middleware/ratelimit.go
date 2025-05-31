package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	redis_storage "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
)

func NewRateLimit(conn redis.UniversalClient, max int, period time.Duration) fiber.Handler {
	if conn == nil {
		panic("redis connection is nil")
	}

	if max <= 0 {
		panic("max is zero")
	}

	if period <= 0 {
		panic("period is zero")
	}

	storage := redis_storage.NewFromConnection(conn)

	return limiter.New(limiter.Config{
		Max: max,
		KeyGenerator: func(c *fiber.Ctx) string {
			return fmt.Sprintf("rl:%s:%s", c.IP(), c.Path())
		},
		Expiration:        period,
		Storage:           storage,
		LimiterMiddleware: limiter.SlidingWindow{},
	})
}
