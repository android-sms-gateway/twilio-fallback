package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Create a Redis client for testing
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a new Fiber app
	app := fiber.New()

	// Create rate limit middleware with a very low limit for testing
	middleware := NewRateLimit(rdb, 2, 10*time.Second)

	// Apply middleware to a test route
	app.Get("/test", middleware, func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Make first request - should succeed
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "2", resp.Header.Get("X-RateLimit-Limit"))
	assert.Equal(t, "1", resp.Header.Get("X-RateLimit-Remaining"))

	// Make second request - should succeed
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "2", resp.Header.Get("X-RateLimit-Limit"))
	assert.Equal(t, "0", resp.Header.Get("X-RateLimit-Remaining"))

	// Make third request - should be rate limited
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	// Check rate limit headers
	assert.Contains(t, resp.Header, "Retry-After")

	fmt.Println(resp.Header)
}
