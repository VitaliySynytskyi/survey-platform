# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy common packages first
COPY pkg/ /app/pkg/

# Copy the API gateway code
COPY api_gateway/ /app/api_gateway/

WORKDIR /app/api_gateway

# Download dependencies
RUN go mod download
RUN go get github.com/google/uuid
RUN go get github.com/hashicorp/consul/api@v1.24.0
RUN go get -d ./...

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/app/main.go

# Final stage
FROM alpine:3.18

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose application port
EXPOSE 8000

# Command to run the application
CMD ["./main"] 