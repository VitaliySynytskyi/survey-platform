package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	AuthServiceURL     string
	SurveyServiceURL   string
	ResponseServiceURL string
	JWTSecret          string
	Port               string
}

func loadConfig() Config {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		SurveyServiceURL:   getEnv("SURVEY_SERVICE_URL", "http://localhost:8082"),
		ResponseServiceURL: getEnv("RESPONSE_SERVICE_URL", "http://localhost:8083"),
		JWTSecret:          getEnv("JWT_SECRET", "your_jwt_secret_key"),
		Port:               getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func setupRouter(config Config) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	// Auth service proxy
	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.Any("/*path", createReverseProxy(config.AuthServiceURL, "/api/v1/auth"))
	}

	// Protected user routes
	userRoutes := r.Group("/api/v1/users")
	{
		userRoutes.Use(jwtAuthMiddleware(config.JWTSecret))
		userRoutes.Any("/*path", createReverseProxy(config.AuthServiceURL, "/api/v1/users"))
	}

	// Protected survey routes
	surveyRoutes := r.Group("/api/v1/surveys")
	{
		// Public route for taking surveys (getting survey details)
		surveyRoutes.GET("/:id", jwtAuthMiddleware(config.JWTSecret), createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))

		// Routes for survey management (CRUD, analytics)
		// Now protected by jwtAuthMiddleware; service layer handles owner/admin logic.
		managedSurveyRoutes := surveyRoutes.Group("")
		managedSurveyRoutes.Use(jwtAuthMiddleware(config.JWTSecret)) // Changed from roleAuthMiddleware("admin")
		{
			managedSurveyRoutes.POST("", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))             // Create survey
			managedSurveyRoutes.PUT("/:id", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))          // Update survey
			managedSurveyRoutes.PATCH("/:id", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))        // Partially update survey
			managedSurveyRoutes.PATCH("/:id/status", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys")) // Update survey status
			managedSurveyRoutes.DELETE("/:id", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))       // Delete survey
			// Survey analytics - service layer will check ownership or admin role.
			// Assuming analytics are now part of survey-service and it checks perms.
			// If analytics were in response-service, it would also need to check X-User-ID/Roles or get survey creator info.
			managedSurveyRoutes.GET("/:id/analytics", createReverseProxy(config.ResponseServiceURL, "/api/v1/surveys")) // Proxy to response-service for analytics
		}

		// Get all surveys for the user (Dashboard) - Authenticated users can see their surveys
		// This remains as jwtAuthMiddleware, service filters by owner or shows all for admin.
		userAccessibleSurveyRoutes := surveyRoutes.Group("")
		userAccessibleSurveyRoutes.Use(jwtAuthMiddleware(config.JWTSecret))
		{
			userAccessibleSurveyRoutes.GET("", createReverseProxy(config.SurveyServiceURL, "/api/v1/surveys"))
		}

		// Routes for survey responses and exports - service layer in response-service will handle owner/admin logic.
		surveyResponseRoutes := surveyRoutes.Group("")
		surveyResponseRoutes.Use(jwtAuthMiddleware(config.JWTSecret)) // Changed from roleAuthMiddleware("admin")
		{
			surveyResponseRoutes.GET("/:id/responses", createReverseProxy(config.ResponseServiceURL, "/api/v1/surveys"))
			surveyResponseRoutes.GET("/:id/responses/export", createReverseProxy(config.ResponseServiceURL, "/api/v1/surveys"))
		}
	}

	// Questions routes - proxied to survey-service. Service layer handles owner/admin logic for survey.
	questionRoutes := r.Group("/api/v1/questions")
	{
		questionRoutes.Use(jwtAuthMiddleware(config.JWTSecret)) // Changed from roleAuthMiddleware("admin")
		questionRoutes.Any("/*path", createReverseProxy(config.SurveyServiceURL, "/api/v1/questions"))
	}

	// Response routes (for submitting new responses) - remains jwtAuthMiddleware
	responseSubmissionRoutes := r.Group("/api/v1/responses")
	{
		// Route for submitting responses - requires authentication
		responseSubmissionRoutes.POST("", jwtAuthMiddleware(config.JWTSecret), createReverseProxy(config.ResponseServiceURL, "/api/v1/responses"))
	}

	return r
}

// createReverseProxy creates a handler to proxy requests to a target service
func createReverseProxy(targetHost string, serviceBasePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client := &http.Client{}

		backendURL, err := url.Parse(targetHost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid backend service host URL"})
			return
		}

		// Construct the path for the backend service.
		// c.Request.URL.Path is the full path from the gateway request, e.g., /api/v1/auth/login
		// serviceBasePath is what the gateway matched, e.g., /api/v1/auth
		// c.Param("path") is the wildcard part, e.g., /login
		// If the backend service expects the full path (e.g. /api/v1/auth/login), we use c.Request.URL.Path.
		// This seems to be the case for auth-service.
		backendURL.Path = c.Request.URL.Path
		backendURL.RawQuery = c.Request.URL.RawQuery

		log.Printf("[API Gateway] Proxying request for %s to %s", c.Request.URL.Path, backendURL.String())

		req, err := http.NewRequest(c.Request.Method, backendURL.String(), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
			return
		}

		// Copy headers from original request
		for name, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Execute the request
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for name, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(name, value)
			}
		}

		// Set status code
		c.Status(resp.StatusCode)

		// Copy response body
		c.Writer.WriteHeader(resp.StatusCode)
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				c.Writer.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
	}
}

// JWT middleware for authentication
func jwtAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := authHeader[7:]

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Check token type
		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type"})
			c.Abort()
			return
		}

		// Store user info in context for downstream services
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])
		c.Set("email", claims["email"])
		c.Set("roles", claims["roles"])

		// Forward the Authorization header to the underlying service
		c.Request.Header.Set("X-User-ID", fmt.Sprintf("%v", claims["user_id"]))
		c.Request.Header.Set("X-User-Roles", fmt.Sprintf("%v", claims["roles"]))

		c.Next()
	}
}

// roleAuthMiddleware checks if the user has one of the required roles
func roleAuthMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawRoles, exists := c.Get("roles")
		if !exists {
			// This should not happen if jwtAuthMiddleware ran successfully
			c.JSON(http.StatusForbidden, gin.H{"error": "roles not found in context"})
			c.Abort()
			return
		}

		userRoles, ok := rawRoles.([]interface{}) // Roles from JWT are often []interface{}
		if !ok {
			// Try asserting to []string if previous assert fails (flexible for different JWT libraries)
			strRoles, okStr := rawRoles.([]string)
			if !okStr {
				c.JSON(http.StatusForbidden, gin.H{"error": "invalid roles format in token"})
				c.Abort()
				return
			}
			// Convert []string to []interface{} to unify logic below, or adapt logic
			userRoles = make([]interface{}, len(strRoles))
			for i, r := range strRoles {
				userRoles[i] = r
			}
		}

		hasRequiredRole := false
		for _, reqRole := range requiredRoles {
			for _, userRole := range userRoles {
				if userRoleStr, ok := userRole.(string); ok && userRoleStr == reqRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	config := loadConfig()
	router := setupRouter(config)

	log.Printf("Starting API Gateway on port %s", config.Port)
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
