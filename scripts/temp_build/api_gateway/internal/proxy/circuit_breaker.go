package proxy

import (
	"errors"
	"sync"
	"time"

	"github.com/VitaliySynytskyi/survey-platform/backend/api_gateway/internal/config"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int

const (
	// CircuitClosed represents a closed circuit (requests flow normally)
	CircuitClosed CircuitState = iota
	// CircuitOpen represents an open circuit (requests are blocked)
	CircuitOpen
	// CircuitHalfOpen represents a half-open circuit (limited requests allowed to test the service)
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mutex           sync.RWMutex
	state           CircuitState
	config          config.CircuitBreakerConfig
	failureCount    uint32
	successCount    uint32
	nextAttempt     time.Time
	lastStateChange time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config config.CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		state:           CircuitClosed,
		config:          config,
		failureCount:    0,
		successCount:    0,
		nextAttempt:     time.Now(),
		lastStateChange: time.Now(),
	}
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	now := time.Now()

	// Check if we need to try moving from open to half-open
	if cb.state == CircuitOpen && now.After(cb.nextAttempt) {
		cb.mutex.RUnlock()
		cb.mutex.Lock()
		defer cb.mutex.Unlock()

		// Change to half-open state
		cb.changeState(CircuitHalfOpen)
		return false
	}

	return cb.state == CircuitOpen
}

// Execute executes a function within the circuit breaker pattern
func (cb *CircuitBreaker) Execute(fn func() error) error {
	// If circuit is open, return error
	if cb.IsOpen() {
		return errors.New("circuit breaker is open")
	}

	// Execute function
	err := fn()

	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	// Update counters based on result
	if err != nil {
		cb.onFailure()
	} else {
		cb.onSuccess()
	}

	return err
}

// onSuccess handles a successful execution
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case CircuitHalfOpen:
		// In half-open state, we need enough consecutive successes to close the circuit
		cb.successCount++
		if cb.successCount >= cb.config.MaxRequests {
			cb.changeState(CircuitClosed)
		}
	case CircuitClosed:
		// In closed state, we reset failure count on success
		cb.failureCount = 0
	}
}

// onFailure handles a failed execution
func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case CircuitClosed:
		// In closed state, we count failures and may trip to open
		cb.failureCount++
		if cb.failureCount >= cb.config.MaxRequests {
			cb.changeState(CircuitOpen)
		}
	case CircuitHalfOpen:
		// In half-open state, a failure immediately trips back to open
		cb.changeState(CircuitOpen)
	}
}

// changeState changes the state of the circuit breaker
func (cb *CircuitBreaker) changeState(newState CircuitState) {
	now := time.Now()
	cb.state = newState
	cb.lastStateChange = now

	switch newState {
	case CircuitOpen:
		// When moving to open, set the retry timeout
		cb.nextAttempt = now.Add(cb.config.Timeout)
		cb.failureCount = 0
		cb.successCount = 0
	case CircuitHalfOpen:
		// When moving to half-open, reset the counters
		cb.failureCount = 0
		cb.successCount = 0
	case CircuitClosed:
		// When moving to closed, reset all counters
		cb.failureCount = 0
		cb.successCount = 0
	}
}
