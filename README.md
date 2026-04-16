# Kinos

Микросервисный интернет-магазин на Go, gRPC, PostgreSQL, Redis и Vue 3.

## Что есть в проекте

- `frontend` - клиентское приложение на Vue 3
- `proto` - protobuf-контракты
- `services/api-service` - HTTP gateway
- `services/user-service` - пользователи, аутентификация, роли
- `services/catalog-service` - категории, производители, товары
- `services/inventory-service` - остатки и склады
- `services/cart-service` - корзина
- `services/order-service` - заказы

## Требования

Для Docker-режима:

- Docker Desktop
- GNU Make

Для локального запуска:

- Go 1.26.x
- Node.js 20+
- Bash

`Makefile` рассчитан на bash. На Windows удобнее использовать Git Bash или WSL.

## Быстрый старт

1. Подготовить env-файлы:

```bash
make env-docker
```

2. Заполнить секреты:

- в `.env` заменить `POSTGRES_PASSWORD`
- в `services/*/.env.docker` заменить все `CHANGE_ME*`

Важно:

- пароль в `services/*/.env.docker` должен совпадать со значением `POSTGRES_PASSWORD` в корневом `.env`

3. Поднять проект:

```bash
make up
```

После старта будут доступны:

- frontend: `http://localhost`
- API gateway: `http://localhost:8080`
- healthcheck: `http://localhost:8080/health`
- user-service gRPC: `localhost:8081`
- catalog-service gRPC: `localhost:8082`
- inventory-service gRPC: `localhost:8083`
- cart-service gRPC: `localhost:8084`
- order-service gRPC: `localhost:8085`

## Основные команды

```bash
make up
make down
make restart
make reset
make test
make test-coverage
make build
make run
make clean
make lint
make deps
make logs
make logs-app
make logs-db
make help
```

## Что делают команды

- `make up` - поднимает весь проект через Docker Compose
- `make down` - останавливает контейнеры
- `make restart` - перезапускает проект
- `make reset` - удаляет контейнеры, volumes и orphan-сервисы
- `make test` - запускает backend-тесты по всем сервисам
- `make test-coverage` - генерирует `coverage.out` и `coverage.html`
- `make build` - собирает все backend-сервисы
- `make run` - запускает frontend и backend локально
- `make clean` - чистит локальные артефакты
- `make lint` - запускает `go vet` и frontend lint
- `make deps` - скачивает Go и frontend зависимости
- `make logs` - показывает все Docker-логи
- `make logs-app` - показывает логи приложений
- `make logs-db` - показывает логи PostgreSQL и Redis

## Локальный запуск

1. Подготовить локальные env-файлы:

```bash
make env-local
```

2. Запустить проект:

```bash
make run
```

Локально будут запущены:

- frontend через Vite
- backend-сервисы как локальные процессы
- PostgreSQL и Redis через Docker Compose

Логи локальных процессов пишутся в `.local/logs`, pid-файлы лежат в `.local/pids`.

Остановить локальные процессы:

```bash
make local-down
```

Посмотреть их логи:

```bash
make local-logs
```

Для локального режима базы доступны на host:

- user-db: `localhost:54321`
- catalog-db: `localhost:54322`
- inventory-db: `localhost:54323`
- order-db: `localhost:54324`

## Лицензия

MIT
