# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git

# Copy the entire backend directory
COPY . .

# Set working directory to the service
WORKDIR /app/liive-message-storer

# Download dependencies
RUN go mod download

# Build with debug information
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o main ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install delve debugger
RUN apk add --no-cache libc6-compat

# Copy binary from builder
COPY --from=builder /app/liive-message-storer/main .

# Expose Message Storer port and debug port
EXPOSE 8084 40004

# Set environment variables
ENV PORT=8084

# Run with or without debugger based on environment variable
CMD ["./main"] 