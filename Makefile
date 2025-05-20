ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_COMPOSE := docker compose
GO_CONTAINER := app
SERVER_CONTAINER := air
DOCKER_COMPOSE_FILE := $(ROOT_DIR)/docker-compose.yaml
ENV_FILE := $(ROOT_DIR)/.env

.PHONY: help setup install build up start down clean logs ps all db-seed db-migrate db-fresh run-inside

.DEFAULT_GOAL := help

setup: ## copy env
	echo "127.0.0.1 bskt.imp" | sudo tee -a /etc/hosts
	@[ -f $(ENV_FILE) ] && echo .env exists || cp .example.env .env
	mkdir -p app/db/migrations

build: ## build
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build --no-cache --progress=plain

up: ## up
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up

start: ## up -d
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

stop: ## stop all services
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop

down: ## down all services
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

sh: ## Enter Golang container sh
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) bash

restart-server: ## Enter Golang container sh
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down $(SERVER_CONTAINER)
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up $(SERVER_CONTAINER) -d

install: ## first time installation
	make setup
	make build

test-with-coverage: ## run tests with coverage
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) go test -v -coverprofile=coverage.out.tmp ./...
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) cat coverage.out.tmp | grep -v "mock_" > coverage.out
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) go tool cover -func coverage.out
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) rm coverage.out.tmp