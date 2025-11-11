# Stage 1: Build the Go application using Go 1.24
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (excluding configs with secrets)
COPY *.go ./
COPY internal/ ./internal/
COPY client/ ./client/
COPY static/ static/
COPY cmd/ ./cmd/

# Build the application binary. CGO_ENABLED=0 ensures static linkage with Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-backend main.go

# Stage 2: Create a minimal production image
FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/todo-backend .

# Create app directory structure
RUN mkdir -p logs

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

# Copy only non-sensitive config files or create a template
COPY --from=builder /app/configs/ ./configs/
# Remove or sanitize the .env file if it was copied
RUN rm -f configs/.env 2>/dev/null || true

# Create a sample environment file template (optional)
RUN echo "# Environment variables will be set at runtime" > configs/.env.template && \
    echo "# Use Docker secrets, environment variables, or external config management" >> configs/.env.template

# Expose the port the Go application listens on
EXPOSE ${HTTP_PORT:-9003}

# Switch to non-root user
USER appuser

# Health check (adjust endpoint based on your app)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:${HTTP_PORT:-9003}/health || exit 1

# Command to run the executable
CMD ["./todo-backend"]