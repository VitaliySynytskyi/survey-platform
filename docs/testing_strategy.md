# Testing Strategy for Survey Application

## 4.1. Purpose and Objectives of Testing

### Main Purpose
The primary purpose of testing is to verify that the developed Survey Platform prototype meets the defined functional and non-functional requirements, and to identify and document any defects.

### Specific Objectives
1. **Functionality Verification**: Test the correct implementation of core functional modules:
   - Authentication and authorization
   - Survey creation and management
   - Survey responses and submission
   - Response collection and analytics

2. **Stability Assessment**: Evaluate the stability of individual microservices and their interactions.

3. **API Correctness**: Verify API endpoints function correctly and return appropriate responses.

4. **User Interface Evaluation**: Perform basic usability testing of the frontend interface.

## 4.2. Testing Levels and Types

### 4.2.1. Unit Testing (Go Microservices)

Our unit testing approach focuses on testing individual components in isolation, particularly in the service and repository layers of each microservice.

#### Auth Service Testing
We've implemented unit tests for the authentication service layer, focusing on key functionality:
- User registration
- Login authentication
- JWT token validation
- Token refresh mechanisms

Example test from `auth-service/internal/service/auth_test.go`:
```go
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
    })
}
```

#### Repository Layer Testing
We've implemented unit tests for the repository layer to verify correct database interactions:

Example from `auth-service/internal/repository/postgres_test.go`:
```go
func TestGetUserByUsername(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %v", err)
    }
    defer db.Close()

    repo := NewPostgresRepository(db)
    ctx := context.Background()

    t.Run("User exists", func(t *testing.T) {
        username := "testuser"
        userID := 1
        
        // Mock database rows and queries
        rows := sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "first_name", "last_name", "is_active", "created_at", "updated_at"}).
            AddRow(userID, username, "test@example.com", "password_hash", "Test", "User", true, time.Now(), time.Now())

        mock.ExpectQuery("SELECT .* FROM users WHERE username = \\$1").
            WithArgs(username).
            WillReturnRows(rows)
            
        // Execute the method
        user, err := repo.GetUserByUsername(ctx, username)

        // Assertions
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, username, user.Username)
    })
}
```

### 4.2.2. Integration Testing

Our integration tests verify the interaction between different microservices and between services and their databases.

#### Auth-Survey Service Integration
We've created tests to verify the end-to-end flow of authentication and survey management:

Example from `integration_tests/auth_survey_integration_test.go`:
```go
func TestAuthAndSurveyIntegration(t *testing.T) {
    // Skip in CI environments
    if os.Getenv("CI") == "true" {
        t.Skip("Skipping integration test in CI environment")
    }

    // Test user registration
    t.Run("Register User", func(t *testing.T) {
        // Test code for user registration
        // ...
    })

    // Test survey creation with auth token
    t.Run("Create Survey", func(t *testing.T) {
        // Test code for creating a survey using auth token
        // ...
    })
    
    // Additional test cases
    // ...
}
```

This test performs the following integration checks:
1. Register a new user account
2. Use the authentication token to create a new survey
3. Retrieve the survey details
4. Update the survey
5. Delete the survey

These integration tests verify the cross-service functionality and token-based authorization is working correctly.

### 4.2.3. API Testing

We've implemented comprehensive API testing using both automated tests and Postman collections.

#### Handlers/API Tests
Example from `response-service/internal/handlers/response_test.go`:
```go
func TestGetResponses(t *testing.T) {
    mockService := new(MockResponseService)
    router := setupTestRouter(mockService)
    
    t.Run("Get responses for survey", func(t *testing.T) {
        // Test data
        surveyID := "survey123"
        responses := []models.Response{
            // Test response data
            // ...
        }
        
        // Set up mock service expectation
        mockService.On("GetResponses", mock.Anything, surveyID).Return(responses, nil).Once()
        
        // Create request and record response
        req, _ := http.NewRequest("GET", "/api/v1/responses/survey/"+surveyID, nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        // Assert correct response
        assert.Equal(t, http.StatusOK, w.Code)
        // More assertions
        // ...
    })
}
```

#### Postman Collection
We've created a Postman collection (`docs/survey-app-api-tests.postman_collection.json`) that includes tests for:
- Authentication endpoints (register, login, token refresh)
- Survey management endpoints (create, read, update, delete)
- Response submission and retrieval endpoints
- Analytics endpoints

Each request includes test scripts to verify the response status, data structure, and business logic correctness.

### 4.2.4. UI/Functional Testing

We've implemented frontend UI tests using Jest and React Testing Library to verify the correct rendering and functionality of key components.

Example from `frontend/src/components/SurveyForm.test.js`:
```javascript
test('successfully submits form with valid data', async () => {
    renderWithRouter(<SurveyForm surveyId="test-survey-123" />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Customer Satisfaction Survey')).toBeInTheDocument();
    });
    
    // Fill form fields
    const ratingInput = screen.getByLabelText('4');
    fireEvent.click(ratingInput);
    
    const textInput = screen.getByPlaceholderText('Your answer');
    fireEvent.change(textInput, { target: { value: 'Great service!' } });
    
    // Submit the form
    const submitButton = screen.getByText('Submit');
    fireEvent.click(submitButton);
    
    // Verify API call and success message
    // ...
});
```

Our UI tests cover scenarios like:
1. Form rendering with different question types
2. Form validation for required fields
3. Successful submission of survey responses
4. Error handling for API failures

### 4.2.5. Non-functional Testing

While comprehensive non-functional testing is beyond the scope of this prototype, we've performed basic evaluations of:

1. **Error Handling**: Testing how the application responds to invalid inputs, authentication failures, and resource not found scenarios.

2. **Response Time**: Basic measurement of API response times for key endpoints.

3. **Cross-browser Compatibility**: Manual verification of frontend functionality in major browsers (Chrome, Firefox).

## 4.3. Testing Environment and Tools

### Environment
- Local development environment using Docker Compose for running all services
- Isolated test databases for unit and integration tests
- CI/CD testing through GitHub Actions workflows

### Tools
- **Unit Testing**: Go standard testing package with testify/assert for assertions and testify/mock for mocking
- **Database Testing**: go-sqlmock for PostgreSQL testing, mongodb-memory-server for MongoDB testing
- **API Testing**: In-code handler tests using httptest and Postman collections for manual/automated API testing
- **UI Testing**: Jest and React Testing Library for frontend component testing
- **Integration Testing**: Custom Go test suite that uses HTTP clients to test cross-service interactions

## 4.4. Test Cases and Results

### 4.4.1. Authentication Module Testing

| Test ID | Description | Expected Result | Actual Result | Status |
|---------|-------------|-----------------|---------------|--------|
| AUTH-1 | User registration with valid data | User created, JWT token returned | As expected | ✅ Pass |
| AUTH-2 | User registration with existing username | Error: "username already exists" | As expected | ✅ Pass |
| AUTH-3 | User registration with existing email | Error: "email already exists" | As expected | ✅ Pass |
| AUTH-4 | Login with valid credentials | JWT token returned | As expected | ✅ Pass |
| AUTH-5 | Login with invalid password | Error: "invalid username or password" | As expected | ✅ Pass |
| AUTH-6 | Login with non-existent user | Error: "invalid username or password" | As expected | ✅ Pass |
| AUTH-7 | Login with inactive account | Error: "account is disabled" | As expected | ✅ Pass |
| AUTH-8 | Token validation with valid token | Claims extracted successfully | As expected | ✅ Pass |
| AUTH-9 | Token validation with invalid token | Error: "invalid token" | As expected | ✅ Pass |
| AUTH-10 | Token refresh with valid refresh token | New tokens returned | As expected | ✅ Pass |
| AUTH-11 | Access protected route with valid token | Access granted | As expected | ✅ Pass |
| AUTH-12 | Access protected route without token | Error: "unauthorized" | As expected | ✅ Pass |

### 4.4.2. Survey Management Module Testing

| Test ID | Description | Expected Result | Actual Result | Status |
|---------|-------------|-----------------|---------------|--------|
| SURV-1 | Create survey with valid data | Survey created, ID returned | As expected | ✅ Pass |
| SURV-2 | Create survey without authentication | Error: "unauthorized" | As expected | ✅ Pass |
| SURV-3 | Get survey by ID | Survey details returned | As expected | ✅ Pass |
| SURV-4 | Get non-existent survey | Error: "survey not found" | As expected | ✅ Pass |
| SURV-5 | Update survey as owner | Survey updated | As expected | ✅ Pass |
| SURV-6 | Update survey as non-owner | Error: "forbidden" | As expected | ✅ Pass |
| SURV-7 | Delete survey as owner | Survey deleted | As expected | ✅ Pass |
| SURV-8 | Get all surveys | List of surveys returned | As expected | ✅ Pass |
| SURV-9 | Create survey with different question types | Survey created with all question types | As expected | ✅ Pass |

### 4.4.3. Response Module Testing

| Test ID | Description | Expected Result | Actual Result | Status |
|---------|-------------|-----------------|---------------|--------|
| RESP-1 | Submit response with valid data | Response saved, ID returned | As expected | ✅ Pass |
| RESP-2 | Submit response without authentication | Error: "unauthorized" | As expected | ✅ Pass |
| RESP-3 | Submit response with missing required answers | Error: "validation failed" | As expected | ✅ Pass |
| RESP-4 | Get responses for a survey | List of responses returned | As expected | ✅ Pass |
| RESP-5 | Get specific response by ID | Response details returned | As expected | ✅ Pass |
| RESP-6 | Get responses for non-existent survey | Empty array returned | As expected | ✅ Pass |
| RESP-7 | Get user's responses | List of user's responses | As expected | ✅ Pass |

### 4.4.4. Analytics Module Testing

| Test ID | Description | Expected Result | Actual Result | Status |
|---------|-------------|-----------------|---------------|--------|
| ANAL-1 | Get survey statistics | Stats with response count and distribution | As expected | ✅ Pass |
| ANAL-2 | Get stats for survey with no responses | Stats with zero counts | As expected | ✅ Pass |
| ANAL-3 | Calculate correct distribution for multiple choice | Accurate percentage distribution | As expected | ✅ Pass |
| ANAL-4 | Calculate rating averages | Correct average rating | As expected | ✅ Pass |
| ANAL-5 | Text response collection | List of text responses | As expected | ✅ Pass |

### 4.4.5. API Gateway Testing

| Test ID | Description | Expected Result | Actual Result | Status |
|---------|-------------|-----------------|---------------|--------|
| GATE-1 | Route auth requests to Auth Service | Correct routing and response | As expected | ✅ Pass |
| GATE-2 | Route survey requests to Survey Service | Correct routing and response | As expected | ✅ Pass |
| GATE-3 | Route response requests to Response Service | Correct routing and response | As expected | ✅ Pass |
| GATE-4 | Token validation and header forwarding | User context headers added | As expected | ✅ Pass |
| GATE-5 | CORS handling | Correct CORS headers in response | As expected | ✅ Pass |
| GATE-6 | Error handling for service unavailability | Appropriate error response | As expected | ✅ Pass |

## 4.5. Identified Issues and Resolutions

During testing, we identified several issues that were subsequently resolved:

### Issue 1: Token Validation Race Condition
**Description**: In high concurrency scenarios, token validation sometimes failed due to race conditions in the token validation middleware.

**Resolution**: Implemented thread-safe token validation with proper synchronization and caching of token validation results.

### Issue 2: Survey Question Validation
**Description**: The survey creation endpoint allowed invalid question configurations (e.g., rating questions without min/max values).

**Resolution**: Added comprehensive validation rules in the survey service to ensure all question types have the required configuration parameters.

### Issue 3: MongoDB Connection Pooling
**Description**: Under load testing, the Response Service showed increasing latency due to connection pool exhaustion.

**Resolution**: Optimized MongoDB connection pooling settings and implemented proper connection management.

### Issue 4: Frontend Form Validation
**Description**: Some question types (particularly multiple-choice questions) weren't properly validated before submission.

**Resolution**: Enhanced client-side validation logic to correctly validate all question types before form submission.

### Issue 5: Cross-Service Error Propagation
**Description**: When one service encountered an error, it wasn't always properly communicated to the client with appropriate status codes.

**Resolution**: Standardized error handling across all services and implemented proper error propagation through the API Gateway.

## Conclusion

The testing strategy implemented for the Survey Application provides comprehensive coverage across all critical components and interactions. Through unit, integration, API, and UI testing, we've verified that the application meets the specified requirements and functions correctly.

The identified issues were addressed, resulting in a more robust and reliable application. Future testing efforts should focus on performance testing under load, security testing, and more extensive cross-browser/cross-device UI testing. 