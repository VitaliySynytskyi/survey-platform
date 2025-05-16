# Architectural Decisions for Survey Platform

This document outlines key architectural decisions made for the Survey Platform microservices application.

## 1. Service Discovery with Consul

### Decision
We've implemented service discovery using HashiCorp Consul.

### Context
In a microservices architecture, services need to discover and communicate with each other dynamically. The number of service instances can change due to scaling, failures, or deployments, making hardcoded configurations impractical.

### Implementation Details
- Each microservice registers itself with Consul on startup
- Services provide health check endpoints that Consul monitors
- Services query Consul to discover other services' endpoints
- Centralized service registry provides a single source of truth about available services
- Service registration and deregistration happens automatically

### Benefits
- Dynamic service discovery eliminates the need for hardcoded service addresses
- Automatic health checks ensure only healthy instances receive traffic
- Enables horizontal scaling without configuration changes
- Provides a centralized view of all services and their health

## 2. Circuit Breaker Pattern

### Decision
We've implemented the Circuit Breaker pattern for handling failures in service-to-service communication.

### Context
In a distributed system, service failures are inevitable. When a service becomes unavailable or experiences high latency, it can cause cascading failures throughout the system as other services wait for responses or exhaust resources.

### Implementation Details
- Circuit breakers are implemented at the API Gateway level
- Each downstream service call is wrapped with a circuit breaker
- Circuit breakers monitor for failures and trip to "open" state after threshold is reached
- In open state, calls fail fast without attempting to reach the downstream service
- After a timeout period, circuit moves to "half-open" state to test if downstream service has recovered
- If test calls succeed, circuit returns to "closed" state

### Benefits
- Prevents cascading failures across the system
- Reduces load on struggling services, allowing them to recover
- Enables faster failure responses, improving user experience
- Provides monitoring and alerting for failing services

## 3. Distributed Tracing

### Decision
We've implemented distributed tracing using Correlation IDs.

### Context
In a microservices architecture, a single request often traverses multiple services. This makes it challenging to trace the request path and diagnose issues that span multiple services.

### Implementation Details
- Each incoming request is assigned a unique Correlation ID
- The ID is propagated through all services in HTTP headers
- For asynchronous communication (RabbitMQ), the ID is included in message properties
- All logs include the Correlation ID
- Centralized logging collects logs from all services and allows filtering by Correlation ID

### Benefits
- Enables end-to-end tracing of requests across multiple services
- Facilitates debugging of complex distributed transactions
- Provides visibility into the performance and behavior of the entire system
- Simplifies troubleshooting by correlating logs from different services

## 4. Data Management Strategy

### Decision
We've implemented a hybrid data storage approach with PostgreSQL for relational data and MongoDB for flexible document storage.

### Context
Different types of data have different storage requirements. Survey platform needs both relational data (users, permissions) and flexible schema data (survey responses, question types).

### Implementation Details
- PostgreSQL is used for:
  - User profiles and authentication data
  - Survey metadata (title, description, ownership, access control)
  - Relationships between entities

- MongoDB is used for:
  - Survey questions and structure (flexible schema for different question types)
  - Survey responses (varied structure based on question types)
  - Analytics data (aggregations and reports)

### Benefits
- Each data store is used for its strengths
- Flexible schema in MongoDB allows for easy evolution of survey designs
- Relational integrity in PostgreSQL ensures consistency for critical data
- Optimization of queries for specific data access patterns

## 5. Asynchronous Communication with RabbitMQ

### Decision
We've implemented asynchronous communication between services using RabbitMQ.

### Context
Some operations don't require immediate responses and can be processed asynchronously. This improves system resilience and scalability.

### Implementation Details
- RabbitMQ used for:
  - Survey response processing
  - Notification delivery
  - Analytics calculations
  - Batch operations

- Implemented patterns:
  - Direct exchanges for targeted message delivery
  - Topic exchanges for category-based routing
  - Dead-letter queues for handling failed messages
  - Message acknowledgments for reliable delivery

### Benefits
- Decouples services for better fault isolation
- Enables back-pressure handling during traffic spikes
- Allows for delayed processing of non-critical operations
- Improves overall system resilience and scalability 