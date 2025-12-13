FROM golang:1.24-alpine AS builder

WORKDIR /todo-backend

COPY go.mod go.sum main.go ./
COPY internal/ ./internal/
COPY "cmd/" "./cmd/"
COPY client/ ./client/
COPY static/ ./static/

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s -extldflags '-static'" -o main .

FROM alpine:3.22 AS production
RUN apk --no-cache add ca-certificates curl tzdata wget \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup
WORKDIR /todo-backend

COPY --from=builder /todo-backend/main /todo-backend/main

# Create app directory structure
RUN mkdir -p logs \
    && chown appuser:appgroup /todo-backend/main \
    && chmod +x /todo-backend/main

USER appuser

ARG HTTP_PORT=9003
EXPOSE ${HTTP_PORT}

CMD ["./main"]