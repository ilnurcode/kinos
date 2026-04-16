SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

SERVICES := services/api-service services/user-service services/catalog-service services/inventory-service services/cart-service services/order-service
CACHE_DIR := $(CURDIR)/.cache
GO_CACHE_DIR := $(CACHE_DIR)/go-build
GO_TMP_DIR := $(CACHE_DIR)/go-tmp
LOCAL_STATE_DIR := .local
LOCAL_PID_DIR := $(LOCAL_STATE_DIR)/pids
LOCAL_LOG_DIR := $(LOCAL_STATE_DIR)/logs
DOCKER_ENV_FILES := .env services/user-service/.env.docker services/catalog-service/.env.docker services/inventory-service/.env.docker services/cart-service/.env.docker services/order-service/.env.docker
LOCAL_ENV_FILES := .env frontend/.env services/user-service/.env services/catalog-service/.env services/inventory-service/.env services/cart-service/.env services/order-service/.env
INFRA_SERVICES := user-db catalog-db inventory-db order-db redis

.PHONY: up down restart reset test test-coverage build run clean lint deps logs logs-app logs-db help env-docker env-local infra-up infra-down local-up local-down local-logs

define ensure_files
missing=0; \
for file in $(1); do \
	if [ ! -f "$$file" ]; then \
		echo "Missing required file: $$file"; \
		missing=1; \
	fi; \
done; \
if [ "$$missing" -ne 0 ]; then \
	exit 1; \
fi
endef

define run_go_for_services
mkdir -p "$(GO_CACHE_DIR)" "$(GO_TMP_DIR)"; \
for svc in $(SERVICES); do \
	echo "==> $$svc"; \
	(cd "$$svc" && GOCACHE="$(GO_CACHE_DIR)" GOTMPDIR="$(GO_TMP_DIR)" $(1)); \
done
endef

define start_local_process
mkdir -p "$(LOCAL_PID_DIR)" "$(LOCAL_LOG_DIR)"; \
pid_file="$(LOCAL_PID_DIR)/$(1).pid"; \
log_file="$(LOCAL_LOG_DIR)/$(1).log"; \
if [ -f "$$pid_file" ] && kill -0 "$$(cat "$$pid_file")" 2>/dev/null; then \
	echo "$(1) already running (pid $$(cat "$$pid_file"))"; \
else \
	echo "Starting $(1)..."; \
	nohup bash -lc 'cd "$(2)" && $(3)' > "$$log_file" 2>&1 & \
	echo $$! > "$$pid_file"; \
	echo "$(1) started. Log: $$log_file"; \
fi
endef

define stop_local_process
pid_file="$(LOCAL_PID_DIR)/$(1).pid"; \
if [ -f "$$pid_file" ]; then \
	pid="$$(cat "$$pid_file")"; \
	if kill -0 "$$pid" 2>/dev/null; then \
		echo "Stopping $(1) (pid $$pid)..."; \
		kill "$$pid" || true; \
		wait "$$pid" 2>/dev/null || true; \
	fi; \
	rm -f "$$pid_file"; \
else \
	echo "No pid file for $(1)"; \
fi
endef

# Подготовка env-файлов для Docker
env-docker:
	@[ -f .env ] || cp .env.example .env
	@[ -f services/user-service/.env.docker ] || cp services/user-service/.env.docker.example services/user-service/.env.docker
	@[ -f services/catalog-service/.env.docker ] || cp services/catalog-service/.env.docker.example services/catalog-service/.env.docker
	@[ -f services/inventory-service/.env.docker ] || cp services/inventory-service/.env.docker.example services/inventory-service/.env.docker
	@[ -f services/cart-service/.env.docker ] || cp services/cart-service/.env.docker.example services/cart-service/.env.docker
	@[ -f services/order-service/.env.docker ] || cp services/order-service/.env.docker.example services/order-service/.env.docker
	@echo "Docker env files are ready"

# Подготовка env-файлов для локального запуска
env-local:
	@[ -f .env ] || cp .env.example .env
	@[ -f frontend/.env ] || cp frontend/.env.example frontend/.env
	@[ -f services/user-service/.env ] || cp services/user-service/.env.example services/user-service/.env
	@[ -f services/catalog-service/.env ] || cp services/catalog-service/.env.example services/catalog-service/.env
	@[ -f services/inventory-service/.env ] || cp services/inventory-service/.env.example services/inventory-service/.env
	@[ -f services/cart-service/.env ] || cp services/cart-service/.env.example services/cart-service/.env
	@[ -f services/order-service/.env ] || cp services/order-service/.env.example services/order-service/.env
	@echo "Local env files are ready"

# Запуск всех сервисов через Docker Compose
up:
	@$(call ensure_files,$(DOCKER_ENV_FILES))
	docker compose up -d --build --remove-orphans
	@echo "Services started. Frontend available at http://localhost, API at http://localhost:8080"

# Остановка всех сервисов
down:
	docker compose down --remove-orphans

# Перезапуск сервисов
restart: down up

# Полный сброс Docker-окружения
reset:
	docker compose down -v --remove-orphans

# Запуск тестов
test:
	@$(call run_go_for_services,go test ./...)

# Запуск тестов с покрытием
test-coverage:
	@mkdir -p "$(CACHE_DIR)"
	@GOCACHE="$(GO_CACHE_DIR)" GOTMPDIR="$(GO_TMP_DIR)" go test ./services/api-service/internal/api/users -coverprofile=coverage.out
	@go tool cover -func=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Сборка backend-сервисов
build:
	@$(call run_go_for_services,go build ./...)

# Локальный запуск проекта
run: env-local
	@$(call ensure_files,$(LOCAL_ENV_FILES))
	@$(MAKE) infra-up
	@$(call start_local_process,user-service,services/user-service,set -a; [ -f ./.env ] && . ./.env || true; set +a; exec go run ./cmd/main.go)
	@$(call start_local_process,catalog-service,services/catalog-service,set -a; [ -f ./.env ] && . ./.env || true; set +a; exec go run ./cmd/main.go)
	@$(call start_local_process,inventory-service,services/inventory-service,set -a; [ -f ./.env ] && . ./.env || true; set +a; exec go run ./cmd/main.go)
	@$(call start_local_process,cart-service,services/cart-service,set -a; [ -f ./.env ] && . ./.env || true; set +a; exec go run ./cmd/main.go)
	@$(call start_local_process,order-service,services/order-service,set -a; [ -f ./.env ] && . ./.env || true; set +a; exec go run ./cmd/main.go)
	@$(call start_local_process,api-service,services/api-service,API_HTTP_PORT=8080 USER_GRPC_ADDR=localhost:8081 CATALOG_GRPC_ADDR=localhost:8082 INVENTORY_GRPC_ADDR=localhost:8083 CART_GRPC_ADDR=localhost:8084 ORDER_GRPC_ADDR=localhost:8085 CORS_ALLOWED_ORIGINS=http://localhost,http://localhost:3000,http://localhost:5173 exec go run ./cmd/main.go)
	@$(call start_local_process,frontend,frontend,set -a; [ -f ./.env ] && . ./.env || true; set +a; [ -d node_modules ] || npm ci; exec npm run dev -- --host 0.0.0.0)
	@echo "Local project is starting. Logs stored in $(LOCAL_LOG_DIR)"

# Очистка
clean:
	go clean
	rm -rf "$(CACHE_DIR)" "$(LOCAL_STATE_DIR)" bin
	rm -f coverage.out coverage.html

# Линтинг
lint:
	@$(call run_go_for_services,go vet ./...)
	@cd frontend && npm run lint

# Установка зависимостей
deps:
	@for svc in $(SERVICES) proto; do \
		echo "==> $$svc"; \
		(cd "$$svc" && go mod download); \
	done
	@cd frontend && npm install

# Показать логи сервисов
logs:
	docker compose logs -f

# Показать логи приложения
logs-app:
	docker compose logs -f api-service frontend user-service catalog-service inventory-service cart-service order-service

# Показать логи БД и Redis
logs-db:
	docker compose logs -f user-db catalog-db inventory-db order-db redis

# Поднять только инфраструктуру для локальной разработки
infra-up:
	@$(call ensure_files,.env)
	docker compose up -d $(INFRA_SERVICES)

# Остановить только инфраструктуру
infra-down:
	docker compose stop $(INFRA_SERVICES)

# Полный локальный запуск
local-up: run

# Остановить локальные процессы
local-down:
	@$(call stop_local_process,frontend)
	@$(call stop_local_process,api-service)
	@$(call stop_local_process,order-service)
	@$(call stop_local_process,cart-service)
	@$(call stop_local_process,inventory-service)
	@$(call stop_local_process,catalog-service)
	@$(call stop_local_process,user-service)
	@docker compose stop $(INFRA_SERVICES) || true

# Показать логи локальных процессов
local-logs:
	@if ls "$(LOCAL_LOG_DIR)"/*.log >/dev/null 2>&1; then tail -f "$(LOCAL_LOG_DIR)"/*.log; else echo "No local logs found"; fi

# Помощь
help:
	@echo "Available commands:"
	@echo "  make env-docker    - Create Docker env files from examples"
	@echo "  make env-local     - Create local env files from examples"
	@echo "  make up            - Start all services with Docker Compose"
	@echo "  make down          - Stop all services"
	@echo "  make restart       - Restart all services"
	@echo "  make reset         - Remove containers, volumes and orphans"
	@echo "  make test          - Run backend tests"
	@echo "  make test-coverage - Generate coverage.out and coverage.html"
	@echo "  make build         - Build all backend services"
	@echo "  make run           - Run frontend and backend locally"
	@echo "  make clean         - Clean local artifacts"
	@echo "  make lint          - Run go vet and frontend lint"
	@echo "  make deps          - Download Go and frontend dependencies"
	@echo "  make logs          - Show Docker Compose logs"
	@echo "  make logs-app      - Show application logs"
	@echo "  make logs-db       - Show database and Redis logs"
	@echo "  make infra-up      - Start only databases and Redis"
	@echo "  make infra-down    - Stop only databases and Redis"
	@echo "  make local-up      - Alias for local run"
	@echo "  make local-down    - Stop locally started processes"
	@echo "  make local-logs    - Tail local process logs"
	@echo "  make help          - Show this help message"
