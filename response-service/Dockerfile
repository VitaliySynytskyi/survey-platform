FROM golang:1.20-alpine AS builder

WORKDIR /app

# Set Go environment variables to potentially bypass proxy issues
ENV GOPROXY=direct
ENV GONOSUMDB=*

# Install git, as some go modules might need it for download
RUN apk add --no-cache git

# We create a temporary go.mod and go.sum to cache dependencies
# This layer will only be rebuilt if go.mod or go.sum changes
COPY go.mod go.sum ./
RUN go mod download

# If using CGO, you might need to install gcc and other build tools
# RUN apk add --no-cache gcc musl-dev

COPY . .

# Generate go.mod and go.sum on the fly and then run go mod tidy
RUN rm -f go.mod go.sum && \
    echo 'module github.com/VitaliySynytskyi/survey-platform/response-service' > go.mod && \
    echo '' >> go.mod && \
    echo 'go 1.20' >> go.mod && \
    go get github.com/gin-gonic/gin@v1.9.0 && \
    go get go.mongodb.org/mongo-driver@v1.13.1 && \
    go get github.com/joho/godotenv@v1.5.1 && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o /response-service ./cmd/app

# Final stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /response-service .
# COPY --from=builder /app/.env .  # Temporarily commented out if .env is not found or not strictly needed here

# Expose port 8082 (or your response service port)
EXPOSE 8082

# Command to run the executable
CMD ["./response-service"] 