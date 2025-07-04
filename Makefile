# Makefile для генератора тестовых данных в MinIO

PROJECT_NAME = minio-fixture-generator

# Docker
DOCKER_COMPOSE = docker compose
DOCKER_BUILD = $(DOCKER_COMPOSE) build
DOCKER_UP = $(DOCKER_COMPOSE) up
DOCKER_DOWN = $(DOCKER_COMPOSE) down

# Targets
.PHONY: help build up down logs clean

help:
	@echo "Доступные команды:"
	@echo "  make build      - собрать Docker-образы"
	@echo "  make up         - запустить MinIO и генератор"
	@echo "  make down       - остановить и удалить все контейнеры"
	@echo "  make logs       - вывести логи всех сервисов"
	@echo "  make clean      - удалить все собранные артефакты и тома"

build:
	$(DOCKER_BUILD)

up:
	$(DOCKER_UP)

down:
	$(DOCKER_DOWN)

logs:
	$(DOCKER_COMPOSE) logs -f

clean:
	$(DOCKER_DOWN) -v --remove-orphans
	docker image prune -f
