package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup test environment
func setupTestRouter() *gin.Engine {
	// Use test mode to suppress gin's debug output during tests
	gin.SetMode(gin.TestMode)

	// Create a test configuration
	config := Config{
		AuthServiceURL:     "http://mock-auth-service",
		SurveyServiceURL:   "http://mock-survey-service",
		ResponseServiceURL: "http://mock-response-service",
		JWTSecret:          "test-secret",
		Port:               "8080",
	}

	// Initialize the router
	return setupRouter(config)
}

func TestHealthEndpoint(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Create a test request
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assert response body
	assert.Nil(t, err)
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "api-gateway", response["service"])
}

func TestCORSMiddleware(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Create a preflight OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(w, req)

	// Assert response status code for OPTIONS request
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Check CORS headers
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

func TestJWTAuthMiddleware(t *testing.T) {
	// Setup test router using our test configuration
	router := setupTestRouter()

	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Create a test request to a protected endpoint without Authorization header
		req := httptest.NewRequest("GET", "/api/v1/surveys/me", nil)
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// Assert error message
		assert.Nil(t, err)
		assert.Equal(t, "authorization header is required", response["error"])
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Create a test request with invalid Authorization format
		req := httptest.NewRequest("GET", "/api/v1/surveys/me", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// Assert error message
		assert.Nil(t, err)
		assert.Equal(t, "invalid authorization header format", response["error"])
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		// Create a test request with invalid JWT token
		req := httptest.NewRequest("GET", "/api/v1/surveys/me", nil)
		req.Header.Set("Authorization", "Bearer invalidtokenformat")
		w := httptest.NewRecorder()

		// Serve the request
		router.ServeHTTP(w, req)

		// Assert response status code
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Parse response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		// Assert error message
		assert.Nil(t, err)
		assert.Equal(t, "invalid token", response["error"])
	})
}

// TestReverseProxyRoutes tests that the API Gateway correctly sets up routes
// We can't test the actual proxying behavior without mock servers,
// but we can verify the routes are registered
func TestReverseProxyRoutes(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Test if routes exist by making requests and checking if the error is
	// about the service being unavailable rather than route not found

	t.Run("Auth Routes", func(t *testing.T) {
		routes := []string{
			"/api/v1/auth/login",
			"/api/v1/auth/register",
		}

		for _, route := range routes {
			req := httptest.NewRequest("POST", route, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Service unavailable or unauthorized, but not 404
			assert.NotEqual(t, http.StatusNotFound, w.Code)
		}
	})

	t.Run("Protected User Routes", func(t *testing.T) {
		routes := []string{
			"/api/v1/users/me",
			"/api/v1/users/me/password",
		}

		for _, route := range routes {
			req := httptest.NewRequest("GET", route, nil)
			req.Header.Set("Authorization", "Bearer invalid")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should be unauthorized, not 404
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Survey Routes", func(t *testing.T) {
		routes := []string{
			"/api/v1/surveys/1",
			"/api/v1/surveys/me",
			"/api/v1/surveys/all",
		}

		for _, route := range routes {
			req := httptest.NewRequest("GET", route, nil)
			req.Header.Set("Authorization", "Bearer invalid")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Should be unauthorized, not 404
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Response Routes", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/responses", nil)
		req.Header.Set("Authorization", "Bearer invalid")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should be unauthorized, not 404
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// Integration tests with mock services would be added here
// These would use httptest.Server to create mock services and test
// the full proxy behavior, but that's beyond the scope of this file
func TestIntegrationWithMockServices(t *testing.T) {
	// This is a placeholder for integration tests
	t.Skip("Integration tests would be implemented here with mock servers")
}
