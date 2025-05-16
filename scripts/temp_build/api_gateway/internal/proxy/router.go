package proxy

import (
	"log"
	"net/http"
	"strings"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/discovery"
)

// RouterType represents the type of router
type RouterType int

const (
	// PathPrefix routes by path prefix
	PathPrefix RouterType = iota
	// Exact matches exact path
	Exact
)

// Route represents a route to a service
type Route struct {
	Type        RouterType
	Path        string
	ServiceName string
	StripPrefix bool
}

// ProxyRouter is a router that proxies requests to different services
type ProxyRouter struct {
	routes    []Route
	proxies   map[string]*ServiceProxy
	config    *config.Config
	discovery discovery.ServiceDiscovery
}

// NewProxyRouter creates a new ProxyRouter
func NewProxyRouter(cfg *config.Config) (*ProxyRouter, error) {
	// Initialize service discovery
	discovery, err := discovery.NewServiceDiscovery(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize service discovery: %v", err)
		// Continue without service discovery
	}

	router := &ProxyRouter{
		routes:    defineRoutes(),
		proxies:   make(map[string]*ServiceProxy),
		config:    cfg,
		discovery: discovery,
	}

	// Initialize proxies for each service
	if err := router.initProxies(); err != nil {
		return nil, err
	}

	return router, nil
}

// ServeHTTP implements the http.Handler interface
func (r *ProxyRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Find the appropriate service for this request path
	for _, route := range r.routes {
		if r.matchRoute(route, req.URL.Path) {
			// Get the service proxy
			proxy, ok := r.proxies[route.ServiceName]
			if !ok {
				http.Error(w, "Service not configured", http.StatusInternalServerError)
				return
			}

			// Modify request path if stripping prefix
			if route.StripPrefix {
				req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
				if req.URL.Path == "" {
					req.URL.Path = "/"
				}
			}

			// Log the request (in production, use a proper logging framework)
			log.Printf("Routing request %s %s to service %s", req.Method, req.URL.Path, route.ServiceName)

			// Forward the request to the appropriate service
			proxy.ServeHTTP(w, req)
			return
		}
	}

	// No matching route found
	http.Error(w, "Not Found", http.StatusNotFound)
}

// matchRoute checks if a request path matches a route
func (r *ProxyRouter) matchRoute(route Route, path string) bool {
	switch route.Type {
	case PathPrefix:
		return strings.HasPrefix(path, route.Path)
	case Exact:
		return path == route.Path
	default:
		return false
	}
}

// initProxies initializes the service proxies
func (r *ProxyRouter) initProxies() error {
	// Create proxy for auth service
	authProxy, err := NewServiceProxy("auth", r.config.Services.Auth, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["auth"] = authProxy

	// Create proxy for user service
	userProxy, err := NewServiceProxy("user", r.config.Services.User, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["user"] = userProxy

	// Create proxy for survey service
	surveyProxy, err := NewServiceProxy("survey", r.config.Services.Survey, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["survey"] = surveyProxy

	// Create proxy for survey taking service
	surveyTakingProxy, err := NewServiceProxy("survey_taking", r.config.Services.SurveyTaking, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["survey_taking"] = surveyTakingProxy

	// Create proxy for response processor service
	responseProcessorProxy, err := NewServiceProxy("response_processor", r.config.Services.ResponseProcessor, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["response_processor"] = responseProcessorProxy

	// Create proxy for analytics service
	analyticsProxy, err := NewServiceProxy("analytics", r.config.Services.Analytics, r.config.CircuitBreaker, r.discovery)
	if err != nil {
		return err
	}
	r.proxies["analytics"] = analyticsProxy

	return nil
}

// defineRoutes defines the routing rules
func defineRoutes() []Route {
	return []Route{
		// Auth service routes
		{
			Type:        PathPrefix,
			Path:        "/api/v1/auth",
			ServiceName: "auth",
			StripPrefix: true,
		},

		// User service routes
		{
			Type:        PathPrefix,
			Path:        "/api/v1/users",
			ServiceName: "user",
			StripPrefix: true,
		},

		// Survey service routes (creation/management)
		{
			Type:        PathPrefix,
			Path:        "/api/v1/surveys",
			ServiceName: "survey",
			StripPrefix: true,
		},

		// Survey taking service routes (for responding to surveys)
		{
			Type:        PathPrefix,
			Path:        "/api/v1/take",
			ServiceName: "survey_taking",
			StripPrefix: true,
		},

		// Analytics service routes
		{
			Type:        PathPrefix,
			Path:        "/api/v1/analytics",
			ServiceName: "analytics",
			StripPrefix: true,
		},

		// Health check
		{
			Type:        Exact,
			Path:        "/health",
			ServiceName: "auth", // Use auth service as a proxy for health checks
			StripPrefix: false,
		},
	}
}
