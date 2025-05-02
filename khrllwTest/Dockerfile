# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build tools
RUN apk add --no-cache git gcc musl-dev

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all files
COPY . .

# Build application (using correct path)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache libc6-compat

# Set working directory
WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

# Copy required files
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/.env ./

# Expose port
EXPOSE 8080

# Run application
CMD ["/app/main"]