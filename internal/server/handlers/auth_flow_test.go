package handlers_test

import (
	"encoding/json"
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
	"github.com/android-sms-gateway/twilio-fallback/internal/proxy"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
	"github.com/android-sms-gateway/twilio-fallback/internal/server/middleware"
	"github.com/android-sms-gateway/twilio-fallback/internal/users"
)

func TestAuthFlow(t *testing.T) {
	logger, _ := zap.NewProduction()
	validator := validator.New()
	authService := auth.NewService(logger, "secret", time.Hour)
	userService := users.NewService(nil, nil, logger)
	proxyService := proxy.NewService(proxy.Config{BaseURL: "localhost"})
	authHandler := handlers.NewAuthHandler(logger, validator, authService, userService)
	dashboardHandler := handlers.NewDashboardHandler(proxyService, logger, validator)

	app := fiber.New()

	// Setup auth routes
	authGroup := app.Group("/auth")
	authHandler.Register(authGroup)

	// Setup protected dashboard route
	protected := app.Group("/dashboard", auth.JWTMiddleware(authService), middleware.UserMiddleware(userService, logger))
	dashboardHandler.Register(protected)

	t.Run("successful auth flow", func(t *testing.T) {
		// Step 1: Login to get token
		loginReqBody := `{"sms_gateway_login": "admin", "sms_gateway_password": "password", "twilio_account_sid": "sid", "twilio_auth_token": "token"}`
		loginReq := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(loginReqBody))
		loginReq.Header.Set("Content-Type", "application/json")

		loginResp, err := app.Test(loginReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, loginResp.StatusCode)

		var loginResponse map[string]interface{}
		require.NoError(t, json.NewDecoder(loginResp.Body).Decode(&loginResponse))
		require.Contains(t, loginResponse, "token")
		token := loginResponse["token"].(string)
		require.NotEmpty(t, token)

		// Step 2: Access protected dashboard with token
		dashboardReq := httptest.NewRequest(http.MethodGet, "/dashboard/", nil)
		dashboardReq.Header.Set("Authorization", "Bearer "+token)

		dashboardResp, err := app.Test(dashboardReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, dashboardResp.StatusCode)

		var dashboardResponse map[string]interface{}
		require.NoError(t, json.NewDecoder(dashboardResp.Body).Decode(&dashboardResponse))
		require.Contains(t, dashboardResponse, "message")
		require.Equal(t, "Welcome to the dashboard", dashboardResponse["message"])
		require.Contains(t, dashboardResponse, "user")
		require.Equal(t, "admin", dashboardResponse["user"])
	})
}
