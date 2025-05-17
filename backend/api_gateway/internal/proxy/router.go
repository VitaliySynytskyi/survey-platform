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
	// CustomMatch allows for custom matching logic, e.g., based on method + path pattern
	CustomMatch
)

// Route represents a route to a service
type Route struct {
	Type         RouterType
	PathPattern  string // For PathPrefix and Exact, this is the path. For CustomMatch, it's a base pattern.
	Method       string // HTTP method for CustomMatch routes (e.g., "POST", "GET", "" for any)
	ServiceName  string
	StripPrefix  bool   // Whether to strip PathPattern from the beginning of the URL path
	TargetPrefix string // Optional prefix to add to the path when forwarding to the service
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
		routes:    defineRoutes(), // Routes are now more structured
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
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	originalPath := req.URL.Path
	var matchedRoute *Route

	for i := range r.routes {
		route := &r.routes[i]
		if r.matchRoute(route, req) {
			matchedRoute = route
			break
		}
	}

	if matchedRoute == nil {
		log.Printf("API_GATEWAY: No matching route for %s %s", req.Method, originalPath)
		http.Error(w, "Not Found on API Gateway", http.StatusNotFound)
		return
	}

	proxy, ok := r.proxies[matchedRoute.ServiceName]
	if !ok {
		log.Printf("API_GATEWAY: Service '%s' not configured for route %s %s", matchedRoute.ServiceName, req.Method, originalPath)
		http.Error(w, "Service not configured", http.StatusInternalServerError)
		return
	}

	newPath := originalPath
	if matchedRoute.StripPrefix {
		newPath = strings.TrimPrefix(originalPath, matchedRoute.PathPattern)
	}

	// Handle specific path transformations for custom routes if StripPrefix was not enough
	// This is where we ensure the service gets the path it expects.
	if matchedRoute.Type == CustomMatch {
		parts := strings.Split(strings.Trim(originalPath, "/"), "/")
		// Example: /api/v1/surveys/{id}/responses
		// parts = ["api", "v1", "surveys", "{id}", "responses"]
		if strings.HasSuffix(matchedRoute.PathPattern, "/:id/responses") && len(parts) >= 5 {
			surveyID := parts[3]
			newPath = "/" + surveyID + "/responses"
		} else if strings.HasSuffix(matchedRoute.PathPattern, "/:id/public") && len(parts) >= 5 {
			// for /api/v1/take/{id}/public - this is handled by PathPrefix + StripPrefix now.
			// but if it were /api/v1/surveys/{id}/public for survey_taking_service
			surveyID := parts[3]
			newPath = "/" + surveyID + "/public"
		} else if strings.HasSuffix(matchedRoute.PathPattern, "/:id/results") && len(parts) >= 5 {
			surveyID := parts[3]
			newPath = "/" + surveyID + "/results"
		}
	}

	if matchedRoute.TargetPrefix != "" {
		newPath = matchedRoute.TargetPrefix + newPath
	}

	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}
	req.URL.Path = newPath

	log.Printf("API_GATEWAY: Routing request original: %s %s, new_path: %s, to service: %s",
		req.Method, originalPath, req.URL.Path, matchedRoute.ServiceName)

	proxy.ServeHTTP(w, req)
}

// matchRoute checks if a request path and method matches a route
func (r *ProxyRouter) matchRoute(route *Route, req *http.Request) bool {
	// Check method for CustomMatch routes
	if route.Method != "" && route.Method != req.Method {
		// If a method is specified in the route and it doesn't match the request method, then it's not a match.
		// This check is important for all route types that might have a method specified.
		// However, for PathPrefix and Exact, we typically don't specify method in Route struct,
		// but if we did, this would handle it.
		// For CustomMatch, this is the primary way to distinguish routes with similar path patterns.
		return false
	}

	switch route.Type {
	case PathPrefix:
		return strings.HasPrefix(req.URL.Path, route.PathPattern)
	case Exact:
		return req.URL.Path == route.PathPattern
	case CustomMatch:
		// For CustomMatch, PathPattern includes parameters like /:id/
		// We need to match the structure and ignore the parameter values themselves during matching.
		requestParts := strings.Split(strings.Trim(req.URL.Path, "/"), "/")
		patternParts := strings.Split(strings.Trim(route.PathPattern, "/"), "/")

		if len(requestParts) != len(patternParts) {
			return false
		}

		for i, patternPart := range patternParts {
			if strings.HasPrefix(patternPart, ":") {
				// This is a parameter in the pattern, it matches any corresponding part in the request path
				continue
			}
			if requestParts[i] != patternPart {
				return false
			}
		}
		// If we've gone through all parts and they matched (or were params), then the route matches.
		return true
	default:
		return false
	}
}

// initProxies initializes the service proxies (ensure service names match those in defineRoutes)
func (r *ProxyRouter) initProxies() error {
	serviceConfigs := map[string]config.ServiceConfig{
		"auth":          r.config.Services.Auth,
		"user":          r.config.Services.User,
		"survey":        r.config.Services.Survey,
		"survey_taking": r.config.Services.SurveyTaking,
		"analytics":     r.config.Services.Analytics,
		// Add other services if they exist in config and are used in routes
	}

	for name, serviceCfg := range serviceConfigs {
		// Check if service is configured directly or via Consul
		if serviceCfg.URL != "" || (r.config.Consul.Enabled && r.config.Consul.UseForSD) {
			proxy, err := NewServiceProxy(name, serviceCfg, r.config.CircuitBreaker, r.discovery)
			if err != nil {
				log.Printf("Warning: Failed to create proxy for service %s: %v. Requests to this service may fail.", name, err)
				// Continue to initialize other proxies
			} else {
				r.proxies[name] = proxy
			}
		} else {
			log.Printf("Info: Service %s is not configured (no URL and service discovery for it is not primary). Proxy not created.", name)
		}
	}
	return nil
}

// defineRoutes defines the routing rules. Order matters: more specific routes should come before general ones.
func defineRoutes() []Route {
	return []Route{
		// AUTH SERVICE
		{
			Type:        PathPrefix,
			PathPattern: "/api/v1/auth",
			ServiceName: "auth",
			StripPrefix: true,
		},
		// USER SERVICE
		{
			Type:        PathPrefix,
			PathPattern: "/api/v1/users",
			ServiceName: "user",
			StripPrefix: true,
		},

		// SURVEY TAKING SERVICE
		// POST /api/v1/surveys/{id}/responses -> survey_taking_service /{id}/responses
		{
			Type:        CustomMatch,
			PathPattern: "/api/v1/surveys/:id/responses", // Full pattern for matching
			Method:      "POST",
			ServiceName: "survey_taking",
			StripPrefix: false, // Path transformation handled in ServeHTTP
		},
		// GET /api/v1/surveys/{id}/public -> survey_taking_service /{id}/public
		// (Used by frontend to fetch survey data for respondent to view and answer)
		{
			Type:        CustomMatch,
			PathPattern: "/api/v1/surveys/:id/public", // Full pattern for matching
			Method:      "GET",
			ServiceName: "survey_taking",
			StripPrefix: false, // Path transformation handled in ServeHTTP by CustomMatch logic
		},
		// GET /api/v1/take/{id}/public -> survey_taking_service /{id}/public (Alternative or specific use?)
		// This route might be redundant if /api/v1/surveys/:id/public is the standard way.
		// If it's used for a different purpose or by a different part of the frontend, keep it.
		// For now, I'm assuming /api/v1/surveys/:id/public is the one that needs to point to survey_taking_service.
		{
			Type:        PathPrefix,
			PathPattern: "/api/v1/take/",
			ServiceName: "survey_taking",
			StripPrefix: true,
		},

		// SURVEY SERVICE
		// GET /api/v1/surveys/{id}/results -> survey_service /{id}/results
		{
			Type:        CustomMatch,
			PathPattern: "/api/v1/surveys/:id/results", // Full pattern
			Method:      "GET",
			ServiceName: "survey",
			StripPrefix: false, // Path transformation handled in ServeHTTP
		},
		// General routes for /api/v1/surveys (e.g., GET /api/v1/surveys/{id}, POST /api/v1/surveys for creation)
		// Must be AFTER specific /api/v1/surveys/.../{action} routes
		{
			Type:        PathPrefix,
			PathPattern: "/api/v1/surveys",
			ServiceName: "survey",
			StripPrefix: true,
			// TargetPrefix: "/", // survey_service expects /{id} or /
		},

		// ANALYTICS SERVICE
		{
			Type:        PathPrefix,
			PathPattern: "/api/v1/analytics",
			ServiceName: "analytics",
			StripPrefix: true,
		},
	}
}
