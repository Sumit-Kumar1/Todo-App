# Stage 1: Build the Go application using Go 1.24
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (excluding configs with secrets)
COPY *.go ./
COPY "cmd/" "./cmd/"
COPY internal/ ./internal/
COPY client/ ./client/
COPY static/ ./static/
# Copy configs but exclude .env file
COPY configs/ ./configs/

# Build the application binary. CGO_ENABLED=0 ensures static linkage with Alpine.
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-backend main.go

# Stage 2: Create a minimal production image
FROM alpine:latest AS production

RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/todo-backend .

# Create app directory structure
RUN mkdir -p logs
RUN mkdir -p configs

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set proper permissions for configs directory
RUN chown -R appuser:appgroup /app/configs && \
    chmod 755 /app/configs

# Expose the port the Go application listens on
EXPOSE ${HTTP_PORT:-9003}

# Switch to non-root user
USER appuser

# Command to run the executable
CMD ["./todo-backend"]
