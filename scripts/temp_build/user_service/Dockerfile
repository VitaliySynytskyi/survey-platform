# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy all files including local dependencies
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:3.18

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata curl

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

# Expose application port
EXPOSE 8080

# Command to run the application
CMD ["./main"]
