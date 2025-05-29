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
	"github.com/android-sms-gateway/twilio-fallback/internal/server/handlers"
)

func TestAuthFlow(t *testing.T) {
	logger, _ := zap.NewProduction()
	validator := validator.New()
	authService := auth.NewAuthService(logger, "secret", time.Hour)
	authHandler := handlers.NewAuthHandler(logger, validator, authService)
	dashboardHandler := handlers.NewDashboardHandler(logger)

	app := fiber.New()

	// Setup auth routes
	authGroup := app.Group("/auth")
	authGroup.Post("/login", authHandler.Login)

	// Setup protected dashboard route
	protected := app.Group("/dashboard", auth.JWTMiddleware(authService))
	protected.Get("/", dashboardHandler.GetDashboard)

	t.Run("successful auth flow", func(t *testing.T) {
		// Step 1: Login to get token
		loginReqBody := `{"username": "admin", "password": "password"}`
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
