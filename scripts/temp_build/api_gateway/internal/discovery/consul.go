package discovery

import (
	"fmt"
	"log"
	"sync"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
	"github.com/VitaliySynytskyi/survey-platform/backend/pkg/consul"
)

// ServiceDiscovery is the interface for service discovery
type ServiceDiscovery interface {
	GetServiceURL(serviceName string) (string, error)
	Close() error
}

// ConsulServiceDiscovery implements ServiceDiscovery using Consul
type ConsulServiceDiscovery struct {
	client      *consul.Client
	cache       map[string]string
	cacheMutex  sync.RWMutex
	serviceList map[string]string // Map of service names to Consul service names
}

// NewConsulServiceDiscovery creates a new Consul service discovery client
func NewConsulServiceDiscovery(cfg *config.Config) (*ConsulServiceDiscovery, error) {
	if !cfg.Consul.Enabled {
		return nil, fmt.Errorf("consul service discovery is disabled")
	}

	client, err := consul.NewClient(cfg.Consul.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	// Initialize service name mapping
	serviceList := map[string]string{
		"auth":               "auth-service",
		"user":               "user-service",
		"survey":             "survey-service",
		"survey_taking":      "survey-taking-service",
		"response_processor": "response-processor-service",
		"analytics":          "analytics-service",
	}

	return &ConsulServiceDiscovery{
		client:      client,
		cache:       make(map[string]string),
		cacheMutex:  sync.RWMutex{},
		serviceList: serviceList,
	}, nil
}

// GetServiceURL returns the URL for a service
func (c *ConsulServiceDiscovery) GetServiceURL(serviceName string) (string, error) {
	// Check if the service is in our mapping
	consulServiceName, ok := c.serviceList[serviceName]
	if !ok {
		return "", fmt.Errorf("unknown service: %s", serviceName)
	}

	// Check cache first
	c.cacheMutex.RLock()
	if url, ok := c.cache[serviceName]; ok {
		c.cacheMutex.RUnlock()
		return url, nil
	}
	c.cacheMutex.RUnlock()

	// Discover the service
	url, err := c.client.GetServiceURL(consulServiceName, "")
	if err != nil {
		return "", fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	// Update cache
	c.cacheMutex.Lock()
	c.cache[serviceName] = url
	c.cacheMutex.Unlock()

	log.Printf("Discovered service %s at %s", serviceName, url)
	return url, nil
}

// Close closes the service discovery client
func (c *ConsulServiceDiscovery) Close() error {
	// Clean up any resources
	return nil
}

// MockServiceDiscovery is a mock implementation for testing or when Consul is disabled
type MockServiceDiscovery struct {
	serviceURLs map[string]string
}

// NewMockServiceDiscovery creates a new mock service discovery client
func NewMockServiceDiscovery(cfg *config.Config) *MockServiceDiscovery {
	return &MockServiceDiscovery{
		serviceURLs: map[string]string{
			"auth":               cfg.Services.Auth.URL,
			"user":               cfg.Services.User.URL,
			"survey":             cfg.Services.Survey.URL,
			"survey_taking":      cfg.Services.SurveyTaking.URL,
			"response_processor": cfg.Services.ResponseProcessor.URL,
			"analytics":          cfg.Services.Analytics.URL,
		},
	}
}

// GetServiceURL returns the URL for a service
func (m *MockServiceDiscovery) GetServiceURL(serviceName string) (string, error) {
	if url, ok := m.serviceURLs[serviceName]; ok {
		return url, nil
	}
	return "", fmt.Errorf("unknown service: %s", serviceName)
}

// Close closes the service discovery client
func (m *MockServiceDiscovery) Close() error {
	return nil
}

// NewServiceDiscovery creates a service discovery client based on configuration
func NewServiceDiscovery(cfg *config.Config) (ServiceDiscovery, error) {
	if cfg.Consul.Enabled && cfg.Consul.UseForSD {
		return NewConsulServiceDiscovery(cfg)
	}
	return NewMockServiceDiscovery(cfg), nil
}
