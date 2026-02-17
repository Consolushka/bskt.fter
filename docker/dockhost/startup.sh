#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Starting application setup..."

# Выполняем миграции один раз при старте
echo "Running database migrations..."
goose up

echo "Running tests..."

go test ./...
echo "Tests completed successfully"

echo "Starting application..."
exec /build
