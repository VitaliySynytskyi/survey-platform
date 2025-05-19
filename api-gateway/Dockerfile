FROM golang:1.20-alpine AS builder

WORKDIR /app

# Create minimal go.mod file
RUN echo 'module github.com/survey-app/api-gateway' > go.mod && \
    echo '' >> go.mod && \
    echo 'go 1.20' >> go.mod

# Copy source code
COPY . .

# Copy go.mod and go.sum files to leverage Docker cache
# COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# Install git for direct go get
RUN apk add --no-cache git

# Reset go.mod and install specific compatible versions
ENV GOPROXY=direct
ENV GOSUMDB=off
RUN rm -f go.mod go.sum && \
    echo 'module github.com/survey-app/api-gateway' > go.mod && \
    echo '' >> go.mod && \
    echo 'go 1.20' >> go.mod && \
    go get github.com/gin-gonic/gin@v1.9.1 && \
    go get github.com/golang-jwt/jwt/v5@v5.0.0 && \
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