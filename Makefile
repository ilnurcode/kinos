# Makefile для проекта Kinos
# Микросервисная архитектура: api-service, user-service, catalog-service, inventory-service

.PHONY: help build run test clean docker-up docker-down docker-restart proto generate migrate frontend

# ==============================================================================
# Справка
# ==============================================================================
help: ## Показать список всех команд
	@echo "Kinos - Микросервисный проект"
	@echo ""
	@echo "Использование:"
	@echo "  make <команда>"
	@echo ""
	@echo "Основные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ==============================================================================
# Frontend (Vue.js)
# ==============================================================================
frontend-install: ## Установить зависимости frontend
	cd frontend && npm install

frontend-dev: ## Запустить frontend в режиме разработки
	cd frontend && npm run dev

frontend-build: ## Собрать frontend для production
	cd frontend && npm run build

frontend-lint: ## Запустить линтер frontend
	cd frontend && npm run lint

frontend: frontend-dev ## Запустить frontend (alias)

# ==============================================================================
# Docker
# ==============================================================================
docker-up: ## Запустить все контейнеры (фон)
	@echo ">>> Запуск Docker контейнеров..."
	docker-compose up -d

docker-up-build: ## Запустить все контейнеры с пересборкой (фон)
	@echo ">>> Запуск Docker контейнеров с пересборкой..."
	docker-compose up -d --build

docker-down: ## Остановить все контейнеры
	@echo ">>> Остановка Docker контейнеров..."
	docker-compose down

docker-restart: ## Перезапустить все контейнеры
	@echo ">>> Перезапуск Docker контейнеров..."
	docker-compose restart

docker-logs: ## Показать логи всех контейнеров
	docker-compose logs -f

docker-logs-frontend: ## Показать логи frontend
	docker-compose logs -f frontend

docker-logs-api: ## Показать логи api-service
	docker-compose logs -f api-service

docker-logs-user: ## Показать логи user-service
	docker-compose logs -f user-service

docker-logs-catalog: ## Показать логи catalog-service
	docker-compose logs -f catalog-service

docker-logs-inventory: ## Показать логи inventory-service
	docker-compose logs -f inventory-service

docker-ps: ## Показать статус контейнеров
	docker-compose ps

docker-clean: ## Очистить Docker кэш
	@echo ">>> Очистка Docker кэша..."
	docker system prune -f

# ==============================================================================
# Сборка Go проектов
# ==============================================================================
build: build-api build-user build-catalog build-inventory ## Собрать все сервисы

build-api: ## Собрать api-service
	@echo ">>> Сборка api-service..."
	cd services/api-service && go build -o ../../bin/api-service ./cmd

build-user: ## Собрать user-service
	@echo ">>> Сборка user-service..."
	cd services/user-service && go build -o ../../bin/user-service ./cmd

build-catalog: ## Собрать catalog-service
	@echo ">>> Сборка catalog-service..."
	cd services/catalog-service && go build -o ../../bin/catalog-service ./cmd

build-inventory: ## Собрать inventory-service
	@echo ">>> Сборка inventory-service..."
	cd services/inventory-service && go build -o ../../bin/inventory-service ./cmd

# ==============================================================================
# Запуск Go проектов (локально)
# ==============================================================================
run: run-api run-user run-catalog run-inventory ## Запустить все сервисы (отдельные терминалы)

run-api: ## Запустить api-service
	@echo ">>> Запуск api-service на :8080"
	cd services/api-service && go run ./cmd/main.go

run-user: ## Запустить user-service
	@echo ">>> Запуск user-service на :8081"
	cd services/user-service && go run ./cmd/main.go

run-catalog: ## Запустить catalog-service
	@echo ">>> Запуск catalog-service на :8082"
	cd services/catalog-service && go run ./cmd/main.go

run-inventory: ## Запустить inventory-service
	@echo ">>> Запуск inventory-service на :8083"
	cd services/inventory-service && go run ./cmd/main.go

# ==============================================================================
# Docker Compose
# ==============================================================================
docker-up: ## Запустить все контейнеры (фон)
	@echo ">>> Запуск Docker контейнеров..."
	docker-compose up -d

docker-up-build: ## Запустить все контейнеры с пересборкой (фон)
	@echo ">>> Запуск Docker контейнеров с пересборкой..."
	docker-compose up -d --build

docker-down: ## Остановить все контейнеры
	@echo ">>> Остановка Docker контейнеров..."
	docker-compose down

docker-restart: ## Перезапустить все контейнеры
	@echo ">>> Перезапуск Docker контейнеров..."
	docker-compose restart

docker-logs: ## Показать логи всех контейнеров
	docker-compose logs -f

docker-logs-api: ## Показать логи api-service
	docker-compose logs -f api-service

docker-logs-user: ## Показать логи user-service
	docker-compose logs -f user-service

docker-logs-catalog: ## Показать логи catalog-service
	docker-compose logs -f catalog-service

docker-logs-inventory: ## Показать логи inventory-service
	docker-compose logs -f inventory-service

docker-ps: ## Показать статус контейнеров
	docker-compose ps

docker-clean: ## Очистить Docker кэш
	@echo ">>> Очистка Docker кэша..."
	docker system prune -f

# ==============================================================================
# Тестирование
# ==============================================================================
test: test-api test-user test-catalog test-inventory ## Запустить тесты всех сервисов

test-api: ## Запустить тесты api-service
	@echo ">>> Тесты api-service..."
	cd services/api-service && go test ./... -v

test-user: ## Запустить тесты user-service
	@echo ">>> Тесты user-service..."
	cd services/user-service && go test ./... -v

test-catalog: ## Запустить тесты catalog-service
	@echo ">>> Тесты catalog-service..."
	cd services/catalog-service && go test ./... -v

test-inventory: ## Запустить тесты inventory-service
	@echo ">>> Тесты inventory-service..."
	cd services/inventory-service && go test ./... -v

# ==============================================================================
# Протоколы (Protobuf)
# ==============================================================================
proto: ## Сгенерировать Go код из .proto файлов
	@echo ">>> Генерация Protobuf..."
	cd proto && go generate ./...

proto-catalog: ## Сгенерировать Protobuf для catalog-service
	@echo ">>> Генерация catalog Protobuf..."
	cd proto/catalog && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		catalog_service.proto

proto-user: ## Сгенерировать Protobuf для user-service
	@echo ">>> Генерация user Protobuf..."
	cd proto/user && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		user_service.proto

proto-inventory: ## Сгенерировать Protobuf для inventory-service
	@echo ">>> Генерация inventory Protobuf..."
	cd proto/inventory && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		inventory_service.proto

# ==============================================================================
# Миграции базы данных
# ==============================================================================
migrate-user-up: ## Применить миграции user-db
	@echo ">>> Миграции user-db (up)..."
	cd services/user-service && go run ./cmd/main.go migrate up

migrate-user-down: ## Откатить миграции user-db
	@echo ">>> Миграции user-db (down)..."
	cd services/user-service && go run ./cmd/main.go migrate down

migrate-catalog-up: ## Применить миграции catalog-db
	@echo ">>> Миграции catalog-db (up)..."
	cd services/catalog-service && go run ./cmd/main.go migrate up

migrate-catalog-down: ## Откатить миграции catalog-db
	@echo ">>> Миграции catalog-db (down)..."
	cd services/catalog-service && go run ./cmd/main.go migrate down

migrate-inventory-up: ## Применить миграции inventory-db
	@echo ">>> Миграции inventory-db (up)..."
	cd services/inventory-service && go run ./cmd/main.go migrate up

migrate-inventory-down: ## Откатить миграции inventory-db
	@echo ">>> Миграции inventory-db (down)..."
	cd services/inventory-service && go run ./cmd/main.go migrate down

# ==============================================================================
# Утилиты
# ==============================================================================
clean: ## Очистить bin и временные файлы
	@echo ">>> Очистка..."
	rm -rf bin/*
	go clean -cache
	go clean -modcache

tidy: ## Обновить зависимости go.mod
	@echo ">>> Обновление зависимостей..."
	cd services/api-service && go mod tidy
	cd services/user-service && go mod tidy
	cd services/catalog-service && go mod tidy
	cd services/inventory-service && go mod tidy
	cd proto && go mod tidy

frontend-tidy: ## Обновить зависимости frontend
	cd frontend && npm install

lint: ## Запустить линтер
	@echo ">>> Линтинг..."
	golangci-lint run ./...

frontend-lint: ## Запустить линтер frontend
	cd frontend && npm run lint

fmt: ## Форматировать код
	@echo ">>> Форматирование..."
	cd services/api-service && go fmt ./...
	cd services/user-service && go fmt ./...
	cd services/catalog-service && go fmt ./...
	cd services/inventory-service && go fmt ./...

vet: ## Запустить go vet
	@echo ">>> go vet..."
	cd services/api-service && go vet ./...
	cd services/user-service && go vet ./...
	cd services/catalog-service && go vet ./...
	cd services/inventory-service && go vet ./...

# ==============================================================================
# Разработка
# ==============================================================================
dev: docker-up ## Запустить среду разработки (Docker)
	@echo ">>> Среда разработки запущена"
	@echo ""
	@echo "Frontend:          http://localhost:80"
	@echo "API Service:       http://localhost:8080"
	@echo "User Service:      localhost:8081 (gRPC)"
	@echo "Catalog Service:   localhost:8082 (gRPC)"
	@echo "Inventory Service: localhost:8083 (gRPC)"
	@echo ""
	@echo "Команды для просмотра логов:"
	@echo "  make docker-logs           - все логи"
	@echo "  make docker-logs-frontend  - frontend логи"
	@echo "  make docker-logs-api       - api-service логи"

dev-stop: docker-down ## Остановить среду разработки
	@echo ">>> Среда разработки остановлена"

# ==============================================================================
# Быстрые команды для разработки
# ==============================================================================
up: docker-up-build ## alias для docker-up-build
down: docker-down ## alias для docker-down
restart: docker-restart ## alias для docker-restart
logs: docker-logs ## alias для docker-logs
ps: docker-ps ## alias для docker-ps

# ==============================================================================
# База данных
# ==============================================================================
db-user-connect: ## Подключиться к user-db
	@echo ">>> Подключение к user-db..."
	docker exec -it user-db psql -U kinos -d user_db

db-catalog-connect: ## Подключиться к catalog-db
	@echo ">>> Подключение к catalog-db..."
	docker exec -it catalog-db psql -U kinos -d catalog_db

db-inventory-connect: ## Подключиться к inventory-db
	@echo ">>> Подключение к inventory-db..."
	docker exec -it inventory-db psql -U kinos -d inventory_db

db-warehouse-connect: ## Подключиться к warehouse-db
	@echo ">>> Подключение к warehouse-db..."
	docker exec -it warehouse-db psql -U kinos -d warehouse_db

db-user-shell: ## Открыть shell в user-db контейнере
	docker exec -it user-db bash

db-catalog-shell: ## Открыть shell в catalog-db контейнере
	docker exec -it catalog-db bash

db-inventory-shell: ## Открыть shell в inventory-db контейнере
	docker exec -it inventory-db bash

db-warehouse-shell: ## Открыть shell в warehouse-db контейнере
	docker exec -it warehouse-db bash

db-inventory-connect: ## Подключиться к inventory-db
	@echo ">>> Подключение к inventory-db..."
	docker exec -it inventory-db psql -U kinos -d inventory_db

db-inventory-shell: ## Открыть shell в inventory-db контейнере
	docker exec -it inventory-db bash
