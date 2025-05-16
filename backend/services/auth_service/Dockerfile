# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory that matches go module path
WORKDIR /src/github.com/VitaliySynytskyi/survey-platform

# Copy the entire project
COPY . .

# Download dependencies
RUN cd backend/services/auth_service && go mod download

# Build the application
RUN cd backend/services/auth_service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.18

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /src/github.com/VitaliySynytskyi/survey-platform/backend/services/auth_service/main .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]