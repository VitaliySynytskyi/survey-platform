# Microservice Survey Platform

A microservice platform for conducting surveys with separation of concerns across multiple services.

## Architecture

The application consists of the following microservices:

- **Auth Service** - Authentication and user management (Golang)
- **Survey Service** - Survey creation, retrieval, and management (Golang)
- **Response Service** - Collection and retrieval of survey responses (Golang)
- **API Gateway** - Single entry point for clients (Golang)
- **Frontend** - User interface for researchers and respondents (Vue.js)

### Data Storage

- PostgreSQL - Stores user data, surveys, questions, etc.
- MongoDB - Stores survey responses
- RabbitMQ - Message broker for future event handling

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Git

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/survey-app.git
   cd survey-app
   ```

2. Start the application:
   ```
   docker-compose up -d
   ```

3. The application will be available at:
   - Frontend: http://localhost:80
   - API Gateway: http://localhost:8080
   - Auth Service: http://localhost:8081
   - Survey Service: http://localhost:8082
   - Response Service: http://localhost:8083
   - RabbitMQ Management: http://localhost:15672 (username: rabbitmq, password: rabbitmq)
   - MongoDB: localhost:27017 (username: mongo, password: mongo)
   - PostgreSQL: localhost:5432 (username: postgres, password: postgres, database: survey_db)

## Development

### Service Structure

Each service follows a similar structure:

```
service-name/
├── cmd/
│   └── app/
│       └── main.go        # Entry point
├── internal/
│   ├── config/            # Configuration
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Data models
│   ├── repository/        # Data access
│   └── service/           # Business logic
├── Dockerfile             # Docker build
└── go.mod                 # Go modules
```

## Features

- User registration and authentication
- JWT-based authorization
- Creation and management of surveys
- Multiple question types
- Collection of survey responses
- Result visualization and export

## API Endpoints

### Auth Service (8081)

- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Log in
- `POST /api/v1/auth/refresh` - Refresh JWT token
- `GET /api/v1/users/me` - Get current user
- `PUT /api/v1/users/me` - Update current user

### Survey Service (8082)

- `POST /api/v1/surveys` - Create a survey
- `GET /api/v1/surveys` - List surveys
- `GET /api/v1/surveys/:id` - Get survey details
- `PUT /api/v1/surveys/:id` - Update survey
- `DELETE /api/v1/surveys/:id` - Delete survey
- `POST /api/v1/surveys/:id/questions` - Add a question
- `PUT /api/v1/questions/:id` - Update a question
- `DELETE /api/v1/questions/:id` - Delete a question

### Response Service (8083)

- `POST /api/v1/responses` - Submit a survey response
- `GET /api/v1/responses` - Get responses for a survey
- `GET /api/v1/responses/summary` - Get response summary

## License

This project is licensed under the MIT License. 