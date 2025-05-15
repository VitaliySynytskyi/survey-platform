# Microservice Platform for Surveys

This project is a microservice-based platform for conducting sociological surveys.

## Technology Stack

- **Frontend**: Vue.js 3 (Composition API)
- **Backend**: Golang
- **Databases**:
  - PostgreSQL (relational data: user profiles, roles, survey metadata)
  - MongoDB (non-relational data: survey structures, responses)
- **Message Broker**: RabbitMQ (for asynchronous communication)
- **Service Discovery**: Consul (planned)

## Microservices Architecture

The platform consists of the following microservices:

- **API Gateway**: Entry point for all client requests, routing to appropriate services
- **Authentication Service**: User authentication and JWT token management
- **User Service**: User profile and account management
- **Survey Service**: Survey creation and management
- **Survey Taking Service**: Handling survey responses
- **Response Processor Service**: Asynchronous processing of survey responses
- **Analytics Service**: Data analysis and reporting

## Project Structure

```
.
├── backend/
│   ├── api_gateway/           # API Gateway
│   └── services/
│       ├── auth_service/      # Authentication and authorization service
│       ├── user_service/      # User management
│       ├── survey_service/    # Survey creation and management
│       ├── survey_taking_service/  # Survey response submission
│       ├── response_processor_service/  # Response processing
│       └── analytics_service/ # Data analysis and reporting
├── frontend/                  # Vue.js client application
├── docker/                    # Docker-related files
├── scripts/                   # Utility scripts
└── docs/                      # Documentation
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.21+
- Node.js 18+
- Git

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/VitaliySynytskyi/survey-platform.git
   cd survey-platform
   ```

2. Initialize the development environment:
   ```
   # For Unix/Linux/macOS
   ./scripts/git-init.sh
   
   # For Windows
   .\scripts\git-init.ps1
   ```

3. Start the services with Docker Compose:
   ```
   docker-compose up -d
   ```

4. Access the services:
   - Frontend: http://localhost:80
   - API Gateway: http://localhost:8000
   - PostgreSQL: localhost:5432
   - MongoDB: localhost:27017
   - RabbitMQ Management UI: http://localhost:15672

## Development

### Backend (Go)

Each microservice follows a standard structure:
- `/cmd`: Entry points for applications
- `/internal`: Private application and library code
- `/pkg`: Public library code
- `/api`: API specifications

### Frontend (Vue.js)

The frontend is a Vue.js 3 application using Composition API.

### Testing

The project includes several types of tests:

- **Unit Tests**: For individual components and functions
- **Integration Tests**: For testing microservices interaction
- **E2E Tests**: End-to-end tests using Cypress (frontend)
- **Contract Tests**: Ensuring compatibility between services

Run the tests using:
```
# Run all tests
go test ./...

# Run specific service tests
cd backend/services/auth_service
go test ./...
```

## Documentation

- **API Documentation**: OpenAPI specifications in `docs/api/`
- **Architecture**: System design and patterns in `docs/architecture.md`

## Contributing

1. Create a feature branch (`git checkout -b feature/xyz`)
2. Make your changes
3. Run tests
4. Commit your changes (`git commit -am 'Add new feature'`)
5. Push the branch (`git push origin feature/xyz`)
6. Create a Pull Request