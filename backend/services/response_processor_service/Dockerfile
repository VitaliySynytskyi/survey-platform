# --- Build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy common packages first
COPY pkg/ /app/pkg/

# Copy the service code
COPY services/response_processor_service/ /app/services/response_processor_service/

WORKDIR /app/services/response_processor_service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/app/main.go

# --- Final stage ---
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

# No need to expose port as this is a worker service

RUN adduser -D -g '' appuser
USER appuser

CMD ["./main"] 