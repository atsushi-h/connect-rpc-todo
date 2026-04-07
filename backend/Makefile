BACKEND_CONTAINER_NAME ?= connect-rpc-todo-backend

MIGRATION_DIR ?= backend/db/migrations
NAME          ?=

ifneq (,$(wildcard ../.env))
include ../.env
export $(shell grep -E '^[A-Za-z0-9_]+=' ../.env | cut -d= -f1)
endif

DOCKER_EXEC := docker exec -i $(BACKEND_CONTAINER_NAME)

# -------------------------------------------------------------------
# App
# -------------------------------------------------------------------

.PHONY: build

build:
	$(DOCKER_EXEC) go build ./backend/...

# -------------------------------------------------------------------
# Docker
# -------------------------------------------------------------------

.PHONY: docker-build docker-up docker-down ssh

docker-build:
	docker compose -f ../compose.yaml build

docker-up:
	docker compose -f ../compose.yaml up -d

docker-down:
	docker compose -f ../compose.yaml down

ssh:
	docker exec -it $(BACKEND_CONTAINER_NAME) sh

# -------------------------------------------------------------------
# Migration
# -------------------------------------------------------------------

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	$(DOCKER_EXEC) sh -c 'migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" up'

migrate-down:
	$(DOCKER_EXEC) sh -c 'migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" down 1'

# 使い方: make migrate-create NAME=add_users
migrate-create:
	@test -n "$(NAME)" || (echo "NAME is required, e.g. make migrate-create NAME=add_users"; exit 1)
	$(DOCKER_EXEC) migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

# -------------------------------------------------------------------
# Tools
# -------------------------------------------------------------------

.PHONY: sqlc

sqlc:
	$(DOCKER_EXEC) sqlc generate -f backend/sqlc.yaml
