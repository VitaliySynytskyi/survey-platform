# Consul Integration Plan for Service Discovery

## Overview

This document outlines the plan for integrating Consul as a service discovery solution in the Survey Platform microservices architecture.

## Goals

1. Implement service discovery to eliminate hardcoded service URLs
2. Provide health checking for all microservices
3. Enable dynamic scaling of services
4. Facilitate service version management

## Implementation Steps

### 1. Add Consul Service to Docker Compose

Update `docker-compose.yml` to add the Consul service:

```yaml
consul:
  image: hashicorp/consul:1.15
  container_name: consul
  ports:
    - "8500:8500"  # UI and API
    - "8600:8600/udp"  # DNS interface
  volumes:
    - consul_data:/consul/data
  environment:
    - CONSUL_BIND_INTERFACE=eth0
    - CONSUL_LOCAL_CONFIG={"datacenter":"dc1"}
  restart: unless-stopped
  healthcheck:
    test: ["CMD", "consul", "members"]
    interval: 10s
    timeout: 5s
    retries: 3
```

### 2. Create Consul Client Package

Create a reusable Consul client package in `backend/pkg/consul`:

```go
// backend/pkg/consul/client.go
package consul

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

// Client represents a Consul client
type Client struct {
	client *consulapi.Client
}

// NewClient creates a new Consul client
func NewClient(address string) (*Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = address

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

// RegisterService registers a service with Consul
func (c *Client) RegisterService(serviceID, serviceName, address string, port int, tags []string, healthCheckURL string) error {
	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
	}

	if healthCheckURL != "" {
		registration.Check = &consulapi.AgentServiceCheck{
			HTTP:                           healthCheckURL,
			Interval:                       "15s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		}
	}

	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Printf("Service '%s' registered with Consul", serviceName)
	return nil
}

// DeregisterService deregisters a service from Consul
func (c *Client) DeregisterService(serviceID string) error {
	err := c.client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	log.Printf("Service ID '%s' deregistered from Consul", serviceID)
	return nil
}

// DiscoverService discovers a service by name
func (c *Client) DiscoverService(serviceName string, tag string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, tag, true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	return services, nil
}

// GetServiceURL returns a URL for a discovered service
func (c *Client) GetServiceURL(serviceName string, tag string) (string, error) {
	services, err := c.DiscoverService(serviceName, tag)
	if err != nil {
		return "", err
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances of service '%s' found", serviceName)
	}

	// Load balancing could be implemented here
	service := services[0]
	return fmt.Sprintf("http://%s:%d", service.Service.Address, service.Service.Port), nil
}
```

### 3. Add Health Check Endpoints to All Microservices

For each microservice, add a health check endpoint:

```go
// Example for auth_service
func setupHealthCheck(router *mux.Router) {
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        // Check dependencies (database, etc.)
        if dbErr := checkDatabaseConnection(); dbErr != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            w.Write([]byte(fmt.Sprintf("Service Unhealthy: %v", dbErr)))
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Service Healthy"))
    }).Methods("GET")
}
```

### 4. Integrate Service Registration in Microservices

Update each microservice to register with Consul on startup and deregister on shutdown:

```go
// Example for auth_service
func main() {
    // ... initialize other components

    // Initialize Consul client
    consulClient, err := consul.NewClient("consul:8500")
    if err != nil {
        log.Fatalf("Failed to create Consul client: %v", err)
    }

    // Register service with Consul
    serviceID := fmt.Sprintf("auth-service-%s", uuid.New().String())
    err = consulClient.RegisterService(
        serviceID,
        "auth-service",
        "auth_service", // container hostname
        8080,           // service port
        []string{"v1", "auth"},
        "http://auth_service:8080/health",
    )
    if err != nil {
        log.Fatalf("Failed to register service with Consul: %v", err)
    }

    // Set up graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        log.Println("Shutting down...")
        
        // Deregister service from Consul
        if err := consulClient.DeregisterService(serviceID); err != nil {
            log.Printf("Failed to deregister service: %v", err)
        }
        
        // ... other cleanup
        os.Exit(0)
    }()

    // ... start HTTP server
}
```

### 5. Update API Gateway to Use Service Discovery

Modify the API Gateway to discover service URLs dynamically:

```go
// Example for API Gateway
func setupAuthRoutes(router *mux.Router, consulClient *consul.Client) {
    router.HandleFunc("/auth/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
        // Discover auth service
        authServiceURL, err := consulClient.GetServiceURL("auth-service", "")
        if err != nil {
            http.Error(w, "Auth service unavailable", http.StatusServiceUnavailable)
            log.Printf("Failed to discover auth service: %v", err)
            return
        }

        // Adjust path
        vars := mux.Vars(r)
        path := vars["path"]
        targetURL := fmt.Sprintf("%s/%s", authServiceURL, path)

        // Proxy request
        proxyRequest(w, r, targetURL)
    })
}

// ... similar for other services
```

### 6. Implement Circuit Breaker with Service Discovery

Update the circuit breaker implementation to use service discovery:

```go
// Example implementation
func CreateCircuitBreaker(consulClient *consul.Client, serviceName string) *gobreaker.CircuitBreaker {
    return gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        serviceName,
        MaxRequests: 5,
        Interval:    10 * time.Second,
        Timeout:     30 * time.Second,
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 10 && failureRatio >= 0.5
        },
        OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
            log.Printf("Circuit '%s' changed from '%s' to '%s'", name, from, to)
        },
    })
}

// Usage in API Gateway
authBreaker := CreateCircuitBreaker(consulClient, "auth-service")
// ... 

// In request handler
result, err := authBreaker.Execute(func() (interface{}, error) {
    serviceURL, err := consulClient.GetServiceURL("auth-service", "")
    if err != nil {
        return nil, err
    }

    // Make the actual HTTP request
    resp, err := http.Get(fmt.Sprintf("%s/%s", serviceURL, path))
    // ... handle response
    return resp, err
})
```

### 7. Update Environment Variables and Configuration

Update service configurations to use Consul:

```yaml
# In docker-compose.yml for each service
environment:
  - CONSUL_ADDR=consul:8500
  - SERVICE_NAME=auth-service
  - SERVICE_PORT=8080
```

### Testing Strategy

1. **Basic Registration/Discovery**: Verify services register correctly with Consul
2. **Health Check**: Verify services report health status correctly
3. **Failure Scenarios**: Test service discovery when services are down
4. **Circuit Breaker**: Test circuit breaker integration with service discovery
5. **Performance**: Measure added latency from service discovery

### Rollout Plan

1. Implement and test Consul integration in development environment
2. Roll out in stages to each microservice
3. Implement in API Gateway last, after all services are integrated with Consul
4. Monitor performance and adjust as needed

## Future Enhancements

1. **Service Mesh**: Consider evolving to a service mesh solution (e.g., Consul Connect)
2. **Configuration Management**: Use Consul for centralized configuration
3. **Dynamic Scaling**: Implement auto-scaling based on service health and load

## Resources

- [Consul Documentation](https://www.consul.io/docs)
- [Go Consul API Client](https://github.com/hashicorp/consul/tree/main/api)
- [Circuit Breaker Pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker)