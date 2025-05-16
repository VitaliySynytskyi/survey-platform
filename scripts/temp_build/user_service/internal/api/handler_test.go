package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/user_service/internal/domain"
	"github.com/google/uuid"
)

// MockUserStore is a mock of UserStore for testing
type MockUserStore struct {
	pingFn func() error
}

func (m *MockUserStore) Ping() error {
	return m.pingFn()
}

// We need to implement all methods of the UserStore interface, but for this test
// we're only interested in the Ping method. The rest can be stubs.
func (m *MockUserStore) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *MockUserStore) UpdateUser(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *MockUserStore) ListUsers(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	return nil, nil
}

func (m *MockUserStore) CountUsers(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *MockUserStore) EmailExists(ctx context.Context, email string, excludeID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *MockUserStore) Close() {
}

func TestHealthCheckHandler_Healthy(t *testing.T) {
	// Create a mock store with a successful Ping
	mockStore := &MockUserStore{
		pingFn: func() error {
			return nil
		},
	}

	// Create the handler with the mock store
	handler := &UserHandler{
		userStore: mockStore,
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlerFunc := http.HandlerFunc(handler.HealthCheckHandler)
	handlerFunc.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var status HealthStatus
	if err := json.Unmarshal(rr.Body.Bytes(), &status); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if status.Status != "healthy" {
		t.Errorf("handler returned wrong status: got %v want %v", status.Status, "healthy")
	}

	if status.Database != "connected" {
		t.Errorf("handler returned wrong database status: got %v want %v", status.Database, "connected")
	}
}

func TestHealthCheckHandler_Unhealthy(t *testing.T) {
	// Create a mock store with a failed Ping
	mockStore := &MockUserStore{
		pingFn: func() error {
			return errors.New("database connection error")
		},
	}

	// Create the handler with the mock store
	handler := &UserHandler{
		userStore: mockStore,
	}

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlerFunc := http.HandlerFunc(handler.HealthCheckHandler)
	handlerFunc.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	// Check the response body
	var status HealthStatus
	if err := json.Unmarshal(rr.Body.Bytes(), &status); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}

	if status.Status != "unhealthy" {
		t.Errorf("handler returned wrong status: got %v want %v", status.Status, "unhealthy")
	}

	if status.Database != "disconnected" {
		t.Errorf("handler returned wrong database status: got %v want %v", status.Database, "disconnected")
	}
}
