package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/survey-app/auth-service/internal/models"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetUserRoles(ctx context.Context, userID int) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRepository) AddUserRole(ctx context.Context, userID, roleID int) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *MockRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func TestRegister(t *testing.T) {
	mockRepo := new(MockRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)
	ctx := context.Background()

	t.Run("Successful registration", func(t *testing.T) {
		req := &models.RegisterRequest{
			Username:  "testuser",
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "Test",
			LastName:  "User",
		}

		// Mock repo responses
		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(nil, errors.New("user not found"))
		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(nil, errors.New("user not found"))
		mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(1, nil)
		mockRepo.On("GetRoleByName", ctx, "user").Return(&models.Role{ID: 1, Name: "user"}, nil)
		mockRepo.On("AddUserRole", ctx, 1, 1).Return(nil)
		mockRepo.On("GetUserByID", ctx, 1).Return(&models.User{
			ID:       1,
			Username: req.Username,
			Email:    req.Email,
			Roles:    []string{"user"},
			IsActive: true,
		}, nil)

		// Call the service
		response, err := authService.Register(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, req.Username, response.User.Username)
		assert.Equal(t, req.Email, response.User.Email)
		assert.Contains(t, response.User.Roles, "user")

		// Verify all mock expectations were met
		mockRepo.AssertExpectations(t)
	})

	t.Run("Username already exists", func(t *testing.T) {
		req := &models.RegisterRequest{
			Username: "existinguser",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Mock existing user
		existingUser := &models.User{
			ID:       1,
			Username: req.Username,
		}

		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(existingUser, nil).Once()

		// Call the service
		response, err := authService.Register(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "username already exists", err.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email already exists", func(t *testing.T) {
		req := &models.RegisterRequest{
			Username: "newuser",
			Email:    "existing@example.com",
			Password: "password123",
		}

		// Mock existing email
		existingUser := &models.User{
			ID:    2,
			Email: req.Email,
		}

		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(nil, errors.New("user not found")).Once()
		mockRepo.On("GetUserByEmail", ctx, req.Email).Return(existingUser, nil).Once()

		// Call the service
		response, err := authService.Register(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "email already exists", err.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockRepository)
	authService := NewAuthService(mockRepo, "test-secret", 24)
	ctx := context.Background()

	t.Run("Successful login", func(t *testing.T) {
		// Hash a test password
		passwordHash, _ := models.HashPassword("password123")

		// Create a mock user
		mockUser := &models.User{
			ID:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: passwordHash,
			Roles:        []string{"user"},
			IsActive:     true,
		}

		req := &models.LoginRequest{
			Username: "testuser",
			Password: "password123",
		}

		// Set up mock repository behavior
		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(mockUser, nil)

		// Call the service
		response, err := authService.Login(ctx, req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		assert.NotEmpty(t, response.RefreshToken)
		assert.Equal(t, mockUser.Username, response.User.Username)
		assert.Equal(t, mockUser.Roles, response.User.Roles)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		req := &models.LoginRequest{
			Username: "nonexistentuser",
			Password: "password123",
		}

		// Set up mock repository behavior
		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(nil, errors.New("user not found"))

		// Call the service
		response, err := authService.Login(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid username or password", err.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid password", func(t *testing.T) {
		// Hash a different password than what will be provided
		passwordHash, _ := models.HashPassword("correctpassword")

		// Create a mock user
		mockUser := &models.User{
			ID:           1,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: passwordHash,
			IsActive:     true,
		}

		req := &models.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword", // This doesn't match the hash
		}

		// Set up mock repository behavior
		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(mockUser, nil)

		// Call the service
		response, err := authService.Login(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid username or password", err.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	t.Run("Inactive account", func(t *testing.T) {
		// Create a mock user with inactive status
		mockUser := &models.User{
			ID:       1,
			Username: "inactiveuser",
			Email:    "inactive@example.com",
			IsActive: false,
		}

		req := &models.LoginRequest{
			Username: "inactiveuser",
			Password: "password123",
		}

		// Set up mock repository behavior
		mockRepo.On("GetUserByUsername", ctx, req.Username).Return(mockUser, nil)

		// Call the service
		response, err := authService.Login(ctx, req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "account is disabled", err.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestValidateToken(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtSecret := "test-secret"
	authService := NewAuthService(mockRepo, jwtSecret, 24)

	t.Run("Valid token", func(t *testing.T) {
		// Create a sample token		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{			"user_id":  1,			"username": "testuser",			"type":     "access",			"exp":      time.Now().Add(24 * time.Hour).Unix(),		})

		tokenString, _ := token.SignedString([]byte(jwtSecret))

		// Validate the token
		claims, err := authService.ValidateToken(tokenString)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, float64(1), claims["user_id"])
		assert.Equal(t, "testuser", claims["username"])
		assert.Equal(t, "access", claims["type"])
	})

	t.Run("Invalid token", func(t *testing.T) {
		// Use a token signed with a different secret
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": 1, "username": "testuser", "type": "access", "exp": time.Now().Add(24 * time.Hour).Unix()})

		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		// Validate the token
		claims, err := authService.ValidateToken(tokenString)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, "invalid token", err.Error())
	})

	t.Run("Incorrect token type", func(t *testing.T) {
		// Create a refresh token instead of access token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": 1,
			"type":    "refresh", // Should be "access"
			"exp":     jwt.NewNumericDate(jwt.NewNumericDate(jwt.Now).Add(24 * time.Hour)).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(jwtSecret))

		// Validate the token
		claims, err := authService.ValidateToken(tokenString)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, "invalid token type", err.Error())
	})
}
