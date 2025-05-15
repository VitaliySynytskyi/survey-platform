# Survey Platform Architecture

## Overview

This document describes the architecture of the "Survey Platform" - a microservices-based platform for creating, managing, and analyzing survey data.

## Infrastructure

The platform uses:
- **Docker & Docker Compose** for containerization and local development
- **Kubernetes** (planned) for production deployment
- **Service Discovery** with Consul (planned)

## Services Architecture

The platform consists of the following microservices:

### Backend Services

1. **API Gateway** (`api_gateway`)
   - Entry point for all client requests
   - Routes to appropriate microservices
   - Handles authentication token validation
   - Implements circuit breaker pattern for fault tolerance
   - Port: 8000

2. **Authentication Service** (`auth_service`)
   - Manages user authentication
   - Issues and validates JWT tokens
   - Handles login, logout, and token refresh
   - Storage: PostgreSQL
   - Port: 8080

3. **User Service** (`user_service`)
   - Manages user profiles and accounts
   - Handles user registration, profile updates
   - User role management
   - Storage: PostgreSQL
   - Port: 8081

4. **Survey Service** (`survey_service`)
   - Manages survey creation and configuration
   - Stores survey metadata in PostgreSQL 
   - Stores survey structure and questions in MongoDB
   - Provides APIs for survey management
   - Port: 8082

5. **Survey Taking Service** (`survey_taking_service`)
   - Handles survey response submission
   - Validates responses against survey configurations
   - Publishes responses to RabbitMQ for async processing
   - Storage: MongoDB
   - Port: 8083

6. **Response Processor Service** (`response_processor_service`)
   - Consumes survey responses from RabbitMQ
   - Processes and validates responses
   - Stores processed responses in MongoDB
   - Implements dead-letter queues for error handling
   - Runs as a background service

7. **Analytics Service** (`analytics_service`)
   - Provides data analysis endpoints
   - Aggregates survey results
   - Generates reports and statistics
   - Storage: MongoDB
   - Port: 8084

### Frontend

- Vue.js 3 application with TypeScript
- Communicates with backend via API Gateway
- Implements client-side form validation
- Responsive design for various devices

## Data Flow

1. **Survey Creation Flow**
   - Admin/creator authenticates via Auth Service
   - Creates survey via Survey Service
   - Survey metadata stored in PostgreSQL
   - Survey structure stored in MongoDB

2. **Survey Taking Flow**
   - Respondent accesses survey
   - Responses submitted to Survey Taking Service
   - Validated responses published to RabbitMQ
   - Response Processor Service consumes message
   - Processed responses stored in MongoDB

3. **Analytics Flow**
   - Admin/creator requests analytics
   - Analytics Service retrieves response data
   - Computes aggregations and statistics
   - Returns analyzed data to client

## Communication Patterns

- **Synchronous Communication**: REST APIs between frontend and backend, and between API Gateway and microservices
- **Asynchronous Communication**: RabbitMQ for event-driven processes, like survey response handling

## Security

- JWT-based authentication
- Role-based access control
- HTTPS for all communications (in production)
- Input validation on both client and server sides

## Database Schema

### PostgreSQL Tables

1. **Users**
   - id (PK)
   - email
   - password_hash
   - name
   - created_at
   - updated_at

2. **Roles**
   - id (PK)
   - name
   - permissions
   - created_at
   - updated_at

3. **User_Roles**
   - user_id (FK)
   - role_id (FK)

4. **Surveys** (Metadata)
   - id (PK)
   - title
   - description
   - creator_id (FK to Users)
   - status
   - public_access
   - start_date
   - end_date
   - created_at
   - updated_at

### MongoDB Collections

1. **survey_templates**
   - _id (ObjectId)
   - survey_id (reference to PostgreSQL)
   - title
   - description
   - questions: [
     - _id
     - question_text
     - question_type
     - required
     - options
     - validation_rules
     - logic_jumps
   ]
   - settings
   - created_at
   - updated_at

2. **survey_responses**
   - _id (ObjectId)
   - survey_id
   - respondent_id (if authenticated)
   - anonymous_id (if anonymous)
   - answers: [
     - question_id
     - value
   ]
   - submitted_at
   - metadata (browser, device, etc.)

3. **survey_results**
   - _id (ObjectId)
   - survey_id
   - aggregated_data
   - last_updated

## Service Discovery (Planned)

Integration with Consul will provide:

1. **Service Registration**: Each microservice registers with Consul at startup
2. **Health Checks**: Regular checks to ensure service availability
3. **Dynamic Configuration**: Update service endpoints without redeployment
4. **DNS Interface**: Simplified service discovery through DNS

## Future Enhancements

1. **Circuit Breaker Implementation**
   - Prevent cascading failures
   - Graceful degradation

2. **Caching Layer**
   - Redis for high-speed data caching
   - Improve response times for frequently accessed data

3. **Monitoring & Logging**
   - Centralized logging
   - Metrics collection
   - Alerting 