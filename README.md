# Kinos — Микросервисный интернет-магазин

[![Go](https://img.shields.io/badge/Go-1.25-blue)](https://go.dev)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

Микросервисное приложение интернет-магазина, написанное на Go с использованием gRPC для межсервисной коммуникации.

## 🏗 Архитектура

```
┌─────────────────┐
│  API Service    │ :8080 (HTTP/REST + Web UI)
│  (Frontend)     │
└────────┬────────┘
         │ gRPC
    ┌────┴────┬──────────┐
    │         │          │
┌───▼────┐  ┌─▼──────┐ ┌─▼──────────┐
│ User   │  │Catalog │ │ Inventory  │
│ Service│  │Service │ │ Service    │
│ :8081  │  │:8082   │ │ :8083      │
└───┬────┘  └───┬────┘ └────┬───────┘
    │           │           │
┌───▼────┐  ┌───▼──────┐ ┌──▼────────┐
│PostgreSQL│ │PostgreSQL│ │PostgreSQL │
│user_db   │ │catalog_db│ │inventory_db│
└──────────┘ └──────────┘ └───────────┘
```

## 📦 Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| **api-service** | 8080 | API-шлюз, HTTP/REST, Web UI |
| **user-service** | 8081 | Управление пользователями, аутентификация |
| **catalog-service** | 8082 | Управление каталогом товаров |
| **inventory-service** | 8083 | Управление запасами товаров, резервирование |

## 🚀 Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.25+ (для локальной разработки)

### Запуск через Docker

```bash
# Запустить все сервисы
docker-compose up -d

# Проверить статус
docker-compose ps
```

Сервисы будут доступны по адресам:
- API Service: http://localhost:8080
- User Service (gRPC): localhost:8081
- Catalog Service (gRPC): localhost:8082
- Inventory Service (gRPC): localhost:8083

### Локальная разработка

```bash
# Установить зависимости
go mod tidy

# Запустить api-service
cd services/api-service && go run ./cmd/main.go
```

## 📋 API Endpoints

### Публичные endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/users/register` | Регистрация |
| POST | `/api/users/login` | Вход |
| GET | `/api/catalog/categories` | Список категорий |
| GET | `/api/catalog/products` | Список товаров |
| GET | `/api/inventory` | Запас товара по ID |
| GET | `/api/inventory/list` | Список запасов |
| POST | `/api/inventory/reserve` | Резервирование товара |
| POST | `/api/inventory/release` | Снятие резервирования |

### Защищённые endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/profile` | Профиль пользователя |
| PUT | `/api/profile` | Обновление профиля |

### Admin endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/admin/users` | Список пользователей |
| PUT | `/api/admin/users/role` | Изменение роли |
| POST/PUT/DELETE | `/api/admin/catalog/*` | Управление каталогом |

## 🛠 Разработка

### Makefile команды

```bash
make help     # Показать все команды
make build    # Собрать все сервисы
make test     # Запустить тесты
make clean    # Очистить временные файлы
```

### Запуск тестов

```bash
# Все тесты
go test ./services/...

# Отдельно по сервисам
go test ./services/user-service/...
go test ./services/catalog-service/...
```

## 📁 Структура проекта

```
kinos/
├── proto/                    # Protobuf определения
├── services/
│   ├── api-service/         # API-шлюз
│   ├── user-service/        # Сервис пользователей
│   ├── catalog-service/     # Сервис каталога
│   └── inventory-service/   # Сервис управления запасами
├── docker-compose.yaml      # Docker конфигурация
├── Makefile                 # Автоматизация задач
└── README.md                # Документация
```

## 🔐 Аутентификация

Проект использует JWT-токены для аутентификации:
- **Access Token** — короткоживущий токен (15 минут)
- **Refresh Token** — долгоживущий токен (30 дней, HttpOnly cookie)

## 📊 Технологии

- **Язык:** Go 1.25
- **RPC:** gRPC
- **Web Framework:** Gin
- **База данных:** PostgreSQL 18
- **Миграции:** golang-migrate
- **Валидация:** go-playground/validator
- **Контейнеризация:** Docker, Docker Compose

## 📝 Лицензия

MIT License — см. файл [LICENSE](LICENSE)
