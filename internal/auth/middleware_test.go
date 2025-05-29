package auth_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
)

func TestJWTMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	authService := auth.NewAuthService(logger, "secret", time.Hour)

	// Generate a valid token
	token, err := authService.GenerateToken("test-user")
	require.NoError(t, err)

	app := fiber.New()
	app.Use(auth.JWTMiddleware(authService))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "protected"})
	})

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "message")
		require.Equal(t, "protected", response["message"])
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "missing or malformed API Key", string(body))
	})

	t.Run("invalid token format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "invalid")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "missing or malformed API Key", string(body))
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, "missing or malformed API Key", string(body))
	})
}
