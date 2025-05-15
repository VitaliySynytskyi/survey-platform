# Microservice Platform for Surveys

This project is a microservice-based platform for conducting sociological surveys.

## Technology Stack

- **Frontend**: Vue.js 3 (Composition API)
- **Backend**: Golang
- **Databases**:
  - PostgreSQL (relational data: user profiles, roles, survey metadata)
  - MongoDB (non-relational data: survey structures, responses)
- **Message Broker**: RabbitMQ (for asynchronous communication)

## Project Structure

```
.
├── backend/
│   └── services/
│       ├── auth_service/       # Authentication and authorization service
│       ├── survey_service/     # Survey creation and management
│       └── user_service/       # User management
├── frontend/                   # Vue.js client application
├── docker/                     # Docker-related files
├── scripts/                    # Utility scripts
└── docs/                       # Documentation
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

## Contributing

1. Create a feature branch (`git checkout -b feature/xyz`)
2. Make your changes
3. Run tests
4. Commit your changes (`git commit -am 'Add new feature'`)
5. Push the branch (`git push origin feature/xyz`)
6. Create a Pull Request