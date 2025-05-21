FROM golang:1.19-alpine

WORKDIR /app

# Install required tools
RUN apk add --no-cache gcc musl-dev git

# Copy go.mod and go.sum files
COPY auth-service/go.mod auth-service/go.sum /app/auth-service/
COPY survey-service/go.mod survey-service/go.sum /app/survey-service/
COPY response-service/go.mod response-service/go.sum /app/response-service/

# Download dependencies for each service
RUN cd /app/auth-service && go mod download
RUN cd /app/survey-service && go mod download
RUN cd /app/response-service && go mod download

# Copy the source code
COPY auth-service /app/auth-service
COPY survey-service /app/survey-service
COPY response-service /app/response-service
COPY integration_tests /app/integration_tests

# Run all unit tests
CMD ["sh", "-c", "cd /app/auth-service && go test -v ./... && cd /app/survey-service && go test -v ./... && cd /app/response-service && go test -v ./..."] 