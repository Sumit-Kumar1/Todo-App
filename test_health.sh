#!/bin/bash

# Test script to verify health endpoint functionality

echo "Testing health endpoint..."

# Set minimal environment variables for testing
export HTTP_PORT=9003
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=todo
export DB_USER=postgres
export DB_PASSWORD=root123
export LOG_LEVEL=DEBUG

echo "Starting backend server on port 9003..."

# Start the backend server in background
cd backend
timeout 10s ./todo-backend &
SERVER_PID=$!

# Wait a moment for server to start
sleep 3

echo "Testing health endpoint at http://localhost:9003/health"

# Test the health endpoint
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9003/health)
echo "Health endpoint response code: $response"

if [ "$response" = "200" ]; then
    echo "✅ Health endpoint is working correctly!"
    curl -s http://localhost:9003/health | jq .
else
    echo "❌ Health endpoint failed with response code: $response"
fi

# Test other health endpoints
echo "Testing readiness endpoint..."
ready_response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9003/ready)
echo "Readiness endpoint response code: $ready_response"

echo "Testing liveness endpoint..."
live_response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:9003/live)
echo "Liveness endpoint response code: $live_response"

# Clean up
kill $SERVER_PID 2>/dev/null
echo "Test completed."
