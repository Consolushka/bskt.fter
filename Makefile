ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
DOCKER_COMPOSE := docker compose
GO_CONTAINER := golang_container
DOCKER_COMPOSE_FILE := $(ROOT_DIR)/docker-compose.yaml
ENV_FILE := $(ROOT_DIR)/.env

.PHONY: help setup install build up start down clean logs ps all db-seed db-migrate db-fresh run-inside

.DEFAULT_GOAL := help

setup: ## Копирование файлов
	@[ -f $(ENV_FILE) ] && echo .env exists || cp .env.example .env

build: ## Запускает сборку образов
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) build

up: ## Запускает docker-compose в текущем shell
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up

start: ## Запускает docker-compose в текущем shell
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up -d

stop: ## Запускает docker-compose в текущем shell
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) stop

down: ## Останавливает контейнеры
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

sh: ## Открыть bash в контейнере nuxtjs
	@$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) exec $(GO_CONTAINER) sh

install: ## Установка проекта
	make setup
	make build