#!/bin/bash


echo "Starting application setup..."

# Выполняем миграции один раз при старте
echo "Running database migrations..."
goose up

/build
