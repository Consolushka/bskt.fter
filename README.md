# IMP

Сервис для сбора и сохранения статистики баскетбольных матчей в PostgreSQL.

## Что делает сервис

- запускает распределенный цикл опроса для каждого активного турнира;
- получает статистику игр из внешних провайдеров (`API_NBA`, `INFOBASKET`, `SPORTOTEKA`);
- сохраняет игры, командную и индивидуальную статистику игроков;
- поддерживает встроенный Rate Limiting для соблюдения лимитов внешних API;
- поддерживает ручной запуск обработки через debug HTTP API.

## Текущие источники данных

- `API_NBA` (с поддержкой Rate Limiting на базе `golang.org/x/time/rate`)
- `INFOBASKET`
- `SPORTOTEKA`

Провайдер `CDN_NBA` присутствует в коде, но не реализован полностью.

## Архитектура

Проект построен по слоям (Hexagonal Architecture):

- `app/internal/core` - доменные модели;
- `app/internal/ports` - интерфейсы (контракты);
- `app/internal/adapters` - реализации репозиториев и провайдеров (GORM, API clients);
- `app/internal/service` - бизнес-логика и оркестрация (`TournamentsOrchestrator`, `Scheduler`);
- `app/internal/infra` - низкоуровневые HTTP-клиенты с поддержкой Rate Limiting;
- `app/database/migrations` - SQL-миграции.

Точки входа:

- `app/cmd/scheduler/main.go` - основной процесс планировщика;
- `app/cmd/debug-server/main.go` - отладочный HTTP-сервер на `:8080`.

## Требования

- Go `1.24.x`
- Docker + Docker Compose (для PostgreSQL и контейнеризации приложения)
- GNU Make
- `golangci-lint` (для `make lint`)

## Переменные окружения

Минимальный набор:

- `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`
- `API_SPORT_API_KEY` (для `API_NBA`)
- `SCHEDULER_POLL_INTERVAL` - интервал опроса каждого турнира (минуты, по умолчанию `30`)
- `SCHEDULER_STAGGER_INTERVAL_MINUTES` - задержка между запусками разных турниров (минуты, по умолчанию `5`)
- `API_NBA_RATE_LIMIT_PER_MINUTE` - лимит запросов для API NBA (по умолчанию `10`)
- `INFOBASKET_RATE_LIMIT_PER_MINUTE` - лимит запросов для Infobasket (по умолчанию `25`)
- `SPORTOTEKA_RATE_LIMIT_PER_MINUTE` - лимит запросов для Sportoteka (по умолчанию `25`)

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

3. Запустить миграции:
```bash
make migrate
```

4. Запустить scheduler:
```bash
make run-scheduler
```

## Развертывание в Docker

Проект содержит `docker/dockhost/Dockerfile`, который собирает приложение и запускает его.
**Важно:** Скрипт запуска `startup.sh` автоматически прогоняет тесты `go test ./...` перед стартом бинарного файла. Если тесты не проходят, приложение не будет запущено.

## Планировщик задач

Планировщик работает по модели **Distributed Workers**:
- При старте получает список всех активных турниров.
- Для каждого турнира запускается независимая горутина-воркер.
- Воркеры запускаются последовательно с задержкой `SCHEDULER_STAGGER_INTERVAL_MINUTES` ("шахматный старт"), чтобы распределить нагрузку на внешние API и базу данных.
- Каждый воркер опрашивает данные своего турнира строго раз в `SCHEDULER_POLL_INTERVAL`.
- Прогресс отслеживается через `poll_watermarks` по `tournament_id`.

## Ручной запуск через debug API

Запуск:
```bash
make run-debug
```

Эндпоинты:
- `GET /health` - проверка работоспособности.
- `GET /process/all?from=YYYY-MM-DD&to=YYYY-MM-DD` - запуск обработки всех активных турниров за период.
- `GET /process/tournament?id=N&from=YYYY-MM-DD&to=YYYY-MM-DD` - запуск обработки конкретного турнира по его ID.

Если параметры `from` и `to` не указаны, обрабатывается период с начала текущих UTC-суток.

## Полезные Make-команды

- `make start` / `make stop` - управление контейнером БД
- `make run-scheduler` - запуск планировщика
- `make run-debug` - запуск debug-сервера
- `make migrate` - применение миграций
- `make lint` - проверка кода линтером
- `make test` - быстрый запуск всех тестов
- `make test-with-coverage` - запуск тестов с генерацией отчета о покрытии (автоматически исключает сгенерированные моки для точности)

## Логирование

Поддерживается вывод в консоль, файл и Telegram. Настройка уровней и путей осуществляется через переменные окружения `LOGGER_*`.
