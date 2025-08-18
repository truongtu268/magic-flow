# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o magic-flow ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S magicflow && \
    adduser -u 1001 -S magicflow -G magicflow

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/magic-flow .

# Copy configuration files
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations

# Create necessary directories
RUN mkdir -p /app/data /app/logs /app/uploads && \
    chown -R magicflow:magicflow /app

# Switch to non-root user
USER magicflow

# Expose ports
EXPOSE 8080 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV CONFIG_PATH=/app/configs/config.yaml
ENV LOG_LEVEL=info

# Run the application
CMD ["./magic-flow"]