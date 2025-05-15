package proxy

import (
	"errors"
	"testing"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
)

func TestCircuitBreaker_NewCircuitBreaker(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		Enabled:       true,
		MaxRequests:   5,
		Interval:      30 * time.Second,
		Timeout:       60 * time.Second,
		TripThreshold: 0.5,
	}

	cb := NewCircuitBreaker(cfg)

	if cb.state != CircuitClosed {
		t.Errorf("Expected initial state to be CircuitClosed, got %v", cb.state)
	}

	if cb.config.MaxRequests != 5 {
		t.Errorf("Expected MaxRequests to be 5, got %v", cb.config.MaxRequests)
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		Enabled:       true,
		MaxRequests:   5,
		Interval:      30 * time.Second,
		Timeout:       60 * time.Second,
		TripThreshold: 0.5,
	}

	cb := NewCircuitBreaker(cfg)

	// Execute a successful function
	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cb.state != CircuitClosed {
		t.Errorf("Expected state to remain CircuitClosed, got %v", cb.state)
	}

	if cb.failureCount != 0 {
		t.Errorf("Expected failureCount to be 0, got %v", cb.failureCount)
	}
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		Enabled:       true,
		MaxRequests:   3, // Set lower for faster testing
		Interval:      30 * time.Second,
		Timeout:       1 * time.Second, // Short timeout for testing
		TripThreshold: 0.5,
	}

	cb := NewCircuitBreaker(cfg)

	// Execute a failing function multiple times
	testErr := errors.New("test error")

	// First failure
	err := cb.Execute(func() error {
		return testErr
	})

	if err != testErr {
		t.Errorf("Expected test error, got %v", err)
	}

	if cb.state != CircuitClosed {
		t.Errorf("Expected state to remain CircuitClosed after first failure, got %v", cb.state)
	}

	if cb.failureCount != 1 {
		t.Errorf("Expected failureCount to be 1, got %v", cb.failureCount)
	}

	// Second failure
	err = cb.Execute(func() error {
		return testErr
	})

	if cb.failureCount != 2 {
		t.Errorf("Expected failureCount to be 2, got %v", cb.failureCount)
	}

	// Third failure - should trip the circuit
	err = cb.Execute(func() error {
		return testErr
	})

	if cb.state != CircuitOpen {
		t.Errorf("Expected state to be CircuitOpen after max failures, got %v", cb.state)
	}

	// After circuit is open, the function should not be executed
	called := false
	err = cb.Execute(func() error {
		called = true
		return nil
	})

	if called {
		t.Error("Expected function not to be called when circuit is open")
	}

	if err == nil {
		t.Error("Expected error when circuit is open, got nil")
	}
}

func TestCircuitBreaker_CircuitHalfOpen(t *testing.T) {
	cfg := config.CircuitBreakerConfig{
		Enabled:       true,
		MaxRequests:   3,
		Interval:      30 * time.Second,
		Timeout:       100 * time.Millisecond, // Very short timeout for testing
		TripThreshold: 0.5,
	}

	cb := NewCircuitBreaker(cfg)

	// Trip the circuit
	for i := 0; i < int(cfg.MaxRequests); i++ {
		cb.Execute(func() error {
			return errors.New("failure")
		})
	}

	if cb.state != CircuitOpen {
		t.Errorf("Expected circuit to be open, got %v", cb.state)
	}

	// Wait for the timeout to allow the circuit to transition to half-open
	time.Sleep(cfg.Timeout + 50*time.Millisecond)

	// The next request should change the state to half-open
	cb.IsOpen() // This should transition to half-open

	if cb.state != CircuitHalfOpen {
		t.Errorf("Expected circuit to be half-open after timeout, got %v", cb.state)
	}

	// Success should keep it in half-open state until enough successes
	cb.Execute(func() error {
		return nil
	})

	if cb.state != CircuitHalfOpen {
		t.Errorf("Expected circuit to remain half-open after one success, got %v", cb.state)
	}

	// Additional successes should close the circuit after reaching MaxRequests
	for i := 0; i < int(cfg.MaxRequests)-1; i++ {
		cb.Execute(func() error {
			return nil
		})
	}

	if cb.state != CircuitClosed {
		t.Errorf("Expected circuit to be closed after enough successes, got %v", cb.state)
	}
}
