#!/bin/bash
set -e

# Change to the repository root directory
cd "$(dirname "$0")/../../.."

# Set the test volume name
export POSTGRES_DATA_VOLUME="postgres_data_test"

# Clean up any existing test environment
echo "Cleaning up test environment..."
docker compose down -v

# Start services with test volume
echo "Starting services..."
docker compose up -d

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 10

# Run the integration tests
echo "Running integration tests..."
cd backend/integration_tests
go test -v ./tests/...

# Capture the test exit code
TEST_EXIT_CODE=$?

# Clean up
echo "Cleaning up..."
cd ../..
docker compose down -v

# Exit with the test exit code
exit $TEST_EXIT_CODE 