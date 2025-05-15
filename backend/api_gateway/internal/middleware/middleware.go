package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"golang.org/x/time/rate"
)

// New creates a new handler with middleware
func New(handler http.Handler) http.Handler {
	return handler
}

// AddLogging adds request logging
func AddLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a custom response writer to capture the status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		handler.ServeHTTP(rw, r)

		// Log the request details
		log.Printf(
			"%s %s %s %d %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}

// AddCORS adds CORS headers
func AddCORS(handler http.Handler, cfg config.CORSConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", joinStringSlice(cfg.AllowedOrigins))
		w.Header().Set("Access-Control-Allow-Methods", joinStringSlice(cfg.AllowedMethods))
		w.Header().Set("Access-Control-Allow-Headers", joinStringSlice(cfg.AllowedHeaders))
		w.Header().Set("Access-Control-Expose-Headers", joinStringSlice(cfg.ExposedHeaders))

		if cfg.AllowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if cfg.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", string(cfg.MaxAge))
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}

// AddRateLimiting adds rate limiting middleware
func AddRateLimiting(handler http.Handler, cfg config.RateLimitingConfig) http.Handler {
	// Create a rate limiter with the configured limit and burst
	limiter := rate.NewLimiter(rate.Limit(cfg.Limit), cfg.Burst)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is allowed
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}

// responseWriter is a custom response writer that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader overrides the WriteHeader method to capture the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Helper functions

// joinStringSlice joins a slice of strings with a comma
func joinStringSlice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}

	// For single item, return it
	if len(slice) == 1 {
		return slice[0]
	}

	// Join multiple items with comma
	result := slice[0]
	for _, s := range slice[1:] {
		result += ", " + s
	}

	return result
}
