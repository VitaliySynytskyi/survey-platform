package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserStore is a mock implementation of the UserStore interface
type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserStore) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserStore) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserStore) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStore) EnsureSchema(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserStore) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockUserStore) Close() {
	m.Called()
}

func TestHealthCheck(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		pingErr        error
		expectedStatus int
	}{
		{
			name:           "Healthy service",
			pingErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Database connection error",
			pingErr:        errors.New("database connection error"),
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock store
			mockStore := new(MockUserStore)
			mockStore.On("Ping", mock.Anything).Return(tc.pingErr)

			// Create a test server
			server := &Server{
				userStore: mockStore,
				handler:   &Handler{}, // We don't need a real handler for this test
			}

			// Create a test request
			req := httptest.NewRequest("GET", "/health", nil)

			// Create a recorder to capture the response
			w := httptest.NewRecorder()

			// Call the health check handler
			server.healthCheckHandler(w, req)

			// Check the response
			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			mockStore.AssertExpectations(t)
		})
	}
}
