# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy the entire backend directory
COPY . .

# Set working directory to the service
WORKDIR /app/liive-migrator

# Download dependencies
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/liive-migrator/main .

CMD ["./main"] 