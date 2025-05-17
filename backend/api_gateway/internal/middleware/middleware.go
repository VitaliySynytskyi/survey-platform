package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/pkg/tracing"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		// Handle preflight requests first and return, as these don't go to the backend.
		if r.Method == http.MethodOptions {
			applyCorsHeaders(w, r, cfg)
			w.WriteHeader(http.StatusOK) // Or http.StatusNoContent, but OK is common.
			return
		}

		// For actual requests, call the next handler (which includes the proxy).
		// This will populate w.Header() from the backend via ReverseProxy.copyHeader (which uses Add()).
		handler.ServeHTTP(w, r)

		// AFTER the backend response headers have been added to 'w',
		// we enforce our CORS policy by using Set(), which will overwrite.
		// This ensures the gateway's CORS policy takes precedence and avoids multiple values.
		applyCorsHeaders(w, r, cfg)
	})
}

// applyCorsHeaders applies the CORS headers based on config and request origin.
// This function should be called to set the final CORS headers on the response.
func applyCorsHeaders(w http.ResponseWriter, r *http.Request, cfg config.CORSConfig) {
	requestOrigin := r.Header.Get("Origin")
	allowedOriginValue := "" // The value to set for Access-Control-Allow-Origin

	// Determine if "*" is an allowed origin
	starAllowed := false
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			starAllowed = true
			break
		}
	}

	// Check if the request's origin is specifically allowed
	specificOriginMatch := false
	if requestOrigin != "" {
		for _, o := range cfg.AllowedOrigins {
			if o == requestOrigin {
				allowedOriginValue = requestOrigin
				specificOriginMatch = true
				break
			}
		}
	}

	if specificOriginMatch {
		w.Header().Set("Access-Control-Allow-Origin", allowedOriginValue)
	} else if starAllowed {
		// If specific origin not matched, but "*" is allowed, use "*"
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	// If neither specific match nor "*", Access-Control-Allow-Origin is not set by this function,
	// effectively clearing any previous value if w.Header().Set was used by backend (unlikely due to Add)
	// or leaving it absent if backend didn't send it.
	// Standard http.Header.Set overwrites, so this behavior is implicitly "clearing" or "setting".

	// Set Vary: Origin if we have allowed origins configured and we might have set ACAO.
	// This tells caches that the response might vary based on the Origin header.
	if len(cfg.AllowedOrigins) > 0 && (specificOriginMatch || starAllowed) {
		w.Header().Add("Vary", "Origin") // Use Add as other parts of the response might also Vary
	}

	if len(cfg.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", joinStringSlice(cfg.AllowedMethods))
	}
	if len(cfg.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", joinStringSlice(cfg.AllowedHeaders))
	}
	if len(cfg.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", joinStringSlice(cfg.ExposedHeaders))
	}
	if cfg.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	if cfg.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
	}
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

// AddCorrelationID adds correlation ID middleware for distributed tracing
func AddCorrelationID(handler http.Handler) http.Handler {
	return tracing.Middleware(handler)
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
	return strings.Join(slice, ", ")
}

// RegisterMiddleware реєструє всі middleware
func RegisterMiddleware(r *chi.Mux, cfg *config.Config) {
	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Correlation ID middleware for distributed tracing
	r.Use(tracing.Middleware)

	// Note: CORS configuration is now handled by the AddCORS middleware function
	// Do not register the chi CORS middleware here to avoid conflicts
}
