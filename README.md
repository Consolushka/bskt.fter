# IMP

Сервис для сбора и сохранения статистики баскетбольных матчей в PostgreSQL.

## Что делает сервис

- запускает фоновые задачи по расписанию (`scheduled_tasks`);
- получает статистику игр из внешних провайдеров;
- сохраняет игры, командную и индивидуальную статистику игроков;
- поддерживает ручной запуск обработки через debug HTTP API.

## Текущие источники данных

- `API_NBA`
- `INFOBASKET`
- `SPORTOTEKA`

Провайдер `CDN_NBA` присутствует в коде, но не реализован полностью.

## Архитектура

Проект построен по слоям:

- `app/internal/core` - доменные модели;
- `app/internal/ports` - интерфейсы (контракты);
- `app/internal/adapters` - реализации репозиториев/провайдеров;
- `app/internal/service` - бизнес-оркестрация (scheduler, processor, persistence);
- `app/internal/infra` - HTTP-клиенты и трансформеры внешних API;
- `app/database/migrations` - SQL-миграции.

Точки входа:

- `app/cmd/scheduler/main.go` - основной воркер планировщика;
- `app/cmd/debug-server/main.go` - отладочный HTTP-сервер на `:8080`.

## Требования

- Go `1.23.x`
- Docker + Docker Compose (только для PostgreSQL)
- GNU Make
- `golangci-lint` (для `make lint`)

## Переменные окружения

Минимальный набор:

- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_PORT`
- `API_SPORT_API_KEY` (для `API_NBA`)

Шаблон: `.example.env`.

## Быстрый старт (локальный Go + Docker DB)

1. Подготовить `.env`:

```bash
cp .example.env .env
```

2. Запустить PostgreSQL:

```bash
make start
```

После запуска:

- `db` доступна на `localhost:5432`.

3. Запустить scheduler локально:

```bash
make run-scheduler
```

4. Или запустить debug API локально:

```bash
make run-debug
```

## Миграции

Миграции лежат в `app/database/migrations`.

Применить миграции:

```bash
make migrate
```

Создать новую миграцию:

```bash
make create-migration name=<migration_name>
```

## Планировщик задач

Типы задач:

- `process_american_tournaments_task`
- `process_not_urgent_european_tournaments_task`
- `process_urgent_european_tournaments_task`

Шедулер:

- читает задачи из `scheduled_tasks`;
- для каждой задачи поднимает отдельный обработчик;
- после выполнения пересчитывает `next_execution_at` и обновляет `last_executed_at`.

## Ручной запуск через debug API

Запуск:

```bash
go run ./app/cmd/debug-server
```

Эндпоинты:

- `GET /process/american?from=YYYY-MM-DD&to=YYYY-MM-DD`
- `GET /process/european-urgent?from=YYYY-MM-DD&to=YYYY-MM-DD`
- `GET /process/european-not-urgent?from=YYYY-MM-DD&to=YYYY-MM-DD`

Пример:

```bash
curl "http://localhost:8080/process/american?from=2026-02-10&to=2026-02-12"
```

## Полезные Make-команды

- `make start` - запуск PostgreSQL в фоне
- `make up` - запуск PostgreSQL в foreground
- `make stop` - остановка PostgreSQL
- `make down` - остановка и удаление PostgreSQL контейнера
- `make run-scheduler` - запуск scheduler локально
- `make run-debug` - запуск debug API локально
- `make test-with-coverage` - тесты и coverage локально
- `make lint` - запуск golangci-lint (включая testifylint)
- `make lint-fix` - автоисправления линтера там, где это возможно

## Тесты

Запуск локально:

```bash
go test ./...
```

Запуск в Docker с coverage:

```bash
make test-with-coverage
```

## Структура БД (основные таблицы)

- `leagues`
- `tournaments`
- `tournament_providers`
- `games`
- `teams`
- `players`
- `game_team_stats`
- `game_team_player_stats`
- `scheduled_tasks`

## Логирование

Фабрика логгеров поддерживает:

- console logger;
- file logger (если включен через env);
- telegram logger (если включен через env).

## Ограничения и заметки

- В `stats_provider` часть методов `GetPlayerBio` еще не реализована для некоторых провайдеров.
- Убедитесь, что `.env` содержит корректные значения для БД и API-ключей перед запуском.
