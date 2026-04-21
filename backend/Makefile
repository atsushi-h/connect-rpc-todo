BACKEND_CONTAINER_NAME ?= connect-rpc-todo-backend

MIGRATION_DIR ?= backend/db/migrations
NAME          ?=
PKG           ?=
FILE          ?=
RUN           ?=
TEST_FLAGS    ?= -count=1

ifneq (,$(wildcard ../.env))
include ../.env
export $(shell grep -E '^[A-Za-z0-9_]+=' ../.env | cut -d= -f1)
endif

DOCKER_EXEC := docker exec -i $(BACKEND_CONTAINER_NAME)

# -------------------------------------------------------------------
# App
# -------------------------------------------------------------------

.PHONY: build test test-pkg test-file test-run

build:
	$(DOCKER_EXEC) go build ./backend/...

# -------------------------------------------------------------------
# Test
# -------------------------------------------------------------------

# 全テストを実行
test:
	$(DOCKER_EXEC) sh -c 'cd backend && go test $(TEST_FLAGS) ./...'

# パッケージ単位でテスト実行。例: make test-pkg PKG=./internal/auth
test-pkg:
	@test -n "$(PKG)" || (echo "PKG is required, e.g. make test-pkg PKG=./internal/auth"; exit 1)
	$(DOCKER_EXEC) sh -c 'cd backend && go test $(TEST_FLAGS) $(PKG)'

# 特定ファイルのテストを実行。例: make test-file FILE=internal/auth/jwt_test.go
test-file:
	@test -n "$(FILE)" || (echo "FILE is required, e.g. make test-file FILE=internal/auth/jwt_test.go"; exit 1)
	@dir=$$(dirname $(FILE)); \
	$(DOCKER_EXEC) sh -c 'cd backend && go test $(TEST_FLAGS) ./$$dir'

# 特定テスト名（regexp）を実行。例: make test-run PKG=./internal/auth RUN=TestSignAndValidateJWT
test-run:
	@test -n "$(PKG)" || (echo "PKG is required, e.g. make test-run PKG=./internal/auth RUN=TestSignAndValidateJWT"; exit 1)
	@test -n "$(RUN)" || (echo "RUN is required, e.g. make test-run PKG=./internal/auth RUN=TestSignAndValidateJWT"; exit 1)
	$(DOCKER_EXEC) sh -c 'cd backend && go test $(TEST_FLAGS) -run $(RUN) $(PKG)'

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

.PHONY: lint fmt sqlc

lint:
	$(DOCKER_EXEC) sh -c 'cd backend && golangci-lint run ./...'

fmt:
	$(DOCKER_EXEC) sh -c 'cd backend && golangci-lint fmt ./...'

sqlc:
	$(DOCKER_EXEC) sqlc generate -f backend/sqlc.yaml
