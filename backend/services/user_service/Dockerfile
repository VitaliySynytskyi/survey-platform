# --- Build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy common packages first
COPY pkg/ /app/pkg/

# Copy the service code
COPY services/user_service/ /app/services/user_service/

WORKDIR /app/services/user_service

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/server/main.go

# --- Final stage ---
FROM alpine:3.18

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

RUN adduser -D -g '' appuser
USER appuser

CMD ["./main"] 