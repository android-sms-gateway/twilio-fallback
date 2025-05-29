package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/android-sms-gateway/twilio-fallback/internal/auth"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
)

func TestAuthHandler_Login(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	validator := validator.New()
	authService := auth.NewAuthService(logger, "secret", time.Hour)
	handler := handlers.NewAuthHandler(logger, validator, authService)

	app := fiber.New()
	app.Post("/login", handler.Login)

	t.Run("successful login", func(t *testing.T) {
		reqBody := `{"username": "admin", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "token")
		require.NotEmpty(t, response["token"])
	})

	t.Run("invalid credentials", func(t *testing.T) {
		reqBody := `{"username": "wrong", "password": "wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
		require.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("empty username", func(t *testing.T) {
		reqBody := `{"username": "", "password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
	})

	t.Run("empty password", func(t *testing.T) {
		reqBody := `{"username": "admin", "password": ""}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
	})

	t.Run("missing username", func(t *testing.T) {
		reqBody := `{"password": "password"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
	})

	t.Run("missing password", func(t *testing.T) {
		reqBody := `{"username": "admin"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
	})

	t.Run("very long credentials", func(t *testing.T) {
		longString := strings.Repeat("a", 1000)
		reqBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, longString, longString)
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
	})

	t.Run("special characters", func(t *testing.T) {
		reqBody := `{"username": "admin<>&'", "password": "p@ssw0rd!"}`
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
		require.Equal(t, "invalid credentials", response["error"])
	})

	t.Run("invalid request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&response))
		require.Contains(t, response, "error")
		require.Equal(t, "invalid request", response["error"])
	})
}
