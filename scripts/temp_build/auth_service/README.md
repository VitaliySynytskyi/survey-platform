# Authentication Service

This service is responsible for user authentication and authorization in the Survey Platform.

## Features

- User registration
- User login with JWT token generation
- Token refresh
- User information retrieval
- Role-based access control

## API Endpoints

### POST /auth/register

Register a new user.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "role": "user",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### POST /auth/login

Login with email and password to get JWT tokens.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### POST /auth/refresh

Refresh access token using a refresh token.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900
}
```

### GET /auth/me

Get information about the current user.

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "role": "user",
  "created_at": "2023-01-01T00:00:00Z"
}
```

## Configuration

The service is configured using environment variables. Copy `env.example` to `.env` and adjust the values as needed:

```
# Server settings
PORT=8080

# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=survey_user
DB_PASSWORD=survey_password
DB_NAME=survey_platform
DB_SSL_MODE=disable

# JWT settings
JWT_SECRET=your-secret-key-change-in-production
JWT_ACCESS_TOKEN_EXP=15m
JWT_REFRESH_TOKEN_EXP=168h
```

## Running the Service

### With Docker

```bash
docker build -t auth-service .
docker run -p 8080:8080 --env-file .env auth-service
```

### Without Docker

```bash
go run cmd/server/main.go
```

## Running Tests

```bash
go test ./...
```

For integration tests that require a database:

```bash
go test -tags=integration ./...
``` 