package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/discovery"
)

// ServiceProxy represents an HTTP proxy to a microservice
type ServiceProxy struct {
	serviceName    string
	config         config.ServiceConfig
	circuitConfig  config.CircuitBreakerConfig
	proxy          *httputil.ReverseProxy
	circuitBreaker *CircuitBreaker
	discovery      discovery.ServiceDiscovery
}

// NewServiceProxy creates a new ServiceProxy
func NewServiceProxy(serviceName string, serviceConfig config.ServiceConfig, circuitConfig config.CircuitBreakerConfig, discovery discovery.ServiceDiscovery) (*ServiceProxy, error) {
	// Get the initial service URL (it will be updated dynamically if using service discovery)
	serviceURL, err := url.Parse(serviceConfig.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid service URL: %w", err)
	}

	// Create circuit breaker if enabled
	var cb *CircuitBreaker
	if circuitConfig.Enabled {
		cb = NewCircuitBreaker(circuitConfig)
	}

	// Create a reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(serviceURL)

	// Override the default director to add service-specific headers or do other manipulations
	defaultDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		// If using service discovery, update the target URL
		if discovery != nil {
			if discoveredURL, err := discovery.GetServiceURL(serviceName); err == nil {
				newURL, _ := url.Parse(discoveredURL)
				req.URL.Scheme = newURL.Scheme
				req.URL.Host = newURL.Host
			}
		}

		defaultDirector(req)
		req.Header.Set("X-Proxy-Time", time.Now().String())
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
	}

	// Set custom timeout for the proxy transport
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: serviceConfig.Timeout,
	}

	// Set custom error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, fmt.Sprintf("Service unavailable: %v", err), http.StatusServiceUnavailable)
	}

	return &ServiceProxy{
		serviceName:    serviceName,
		config:         serviceConfig,
		circuitConfig:  circuitConfig,
		proxy:          proxy,
		circuitBreaker: cb,
		discovery:      discovery,
	}, nil
}

// ServeHTTP implements the http.Handler interface
func (s *ServiceProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if circuit breaker is enabled and open
	if s.circuitBreaker != nil && s.circuitBreaker.IsOpen() {
		http.Error(w, "Service temporarily unavailable (circuit open)", http.StatusServiceUnavailable)
		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), s.config.Timeout)
	defer cancel()

	// Use the context-aware request
	r = r.WithContext(ctx)

	// Create a response recorder to capture the response
	recorder := NewResponseRecorder(w)

	// Use circuit breaker if enabled
	if s.circuitBreaker != nil {
		s.circuitBreaker.Execute(func() error {
			s.proxy.ServeHTTP(recorder, r)
			if recorder.Status >= 500 {
				return errors.New("service error")
			}
			return nil
		})
	} else {
		// Directly serve the request through the proxy if no circuit breaker
		s.proxy.ServeHTTP(recorder, r)
	}
}

// ResponseRecorder is a custom response writer that captures the status code
type ResponseRecorder struct {
	http.ResponseWriter
	Status int
}

// NewResponseRecorder creates a new ResponseRecorder
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Status:         http.StatusOK, // Default status
	}
}

// WriteHeader overrides the WriteHeader method to capture the status code
func (r *ResponseRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
