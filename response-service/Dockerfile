FROM golang:1.20-alpine AS builder

WORKDIR /app

# Create minimal go.mod file
RUN echo 'module github.com/survey-app/response-service' > go.mod && \
    echo '' >> go.mod && \
    echo 'go 1.20' >> go.mod

# Copy source code
COPY . .

# Reset go.mod and install specific compatible versions
RUN rm -f go.mod go.sum && \
    echo 'module github.com/survey-app/response-service' > go.mod && \
    echo '' >> go.mod && \
    echo 'go 1.20' >> go.mod && \
    go get github.com/gin-gonic/gin@v1.9.0 && \
    go get go.mongodb.org/mongo-driver@v1.13.1 && \
    go get github.com/joho/godotenv@v1.5.1 && \
    go mod tidy

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Use minimal alpine image
FROM alpine:3.17

WORKDIR /app

# Copy built binary from builder
COPY --from=builder /app/main .

# Run the executable
CMD ["./main"] 