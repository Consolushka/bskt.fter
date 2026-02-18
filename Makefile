ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_COMPOSE := docker compose
DOCKER_COMPOSE_FILE := $(ROOT_DIR)/docker-compose.yaml
ENV_FILE := $(ROOT_DIR)/.env
GOOSE := goose

ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

.PHONY: setup db-up db-start db-stop db-down up start stop down run-scheduler run-debug test-with-coverage migrate create-migration lint lint-fix

.DEFAULT_GOAL := db-up

setup: ## copy env
	@[ -f $(ENV_FILE) ] && echo ".env exists" || cp .example.env .env

db-up: ## run postgres in foreground
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up db

db-start: ## run postgres in background
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d db

db-stop: ## stop postgres container
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop db

db-down: ## stop and remove postgres container
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down db

# Backward-compatible aliases
up: db-up
start: db-start
stop: db-stop
down: db-down

run-scheduler: ## run scheduler locally
	@go run ./app/cmd/scheduler

run-debug: ## run debug server locally
	@go run ./app/cmd/debug-server

test:
	@go test ./...

test-verbose:
	@go test -v ./...

test-with-coverage: ## run tests with coverage locally, excluding mocks
	@go test -coverprofile=coverage.out ./...
	@grep -v "mock_" coverage.out > coverage.out.tmp || true
	@mv coverage.out.tmp coverage.out
	@go tool cover -func coverage.out

lint: ## run golangci-lint (includes testifylint)
	@golangci-lint run

lint-fix: ## run golangci-lint with auto-fix where available
	@golangci-lint run --fix

migrate: ## run goose up locally
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) $(GOOSE) up

create-migration: ## create goose migration locally: make create-migration name=add_some_column
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) $(GOOSE) create $(name) sql
