# CLAUDE.md

Технический контекст проекта **IMP** для работы в этом репозитории.
Роль, стиль общения и инженерные принципы заданы глобально в `~/.claude/CLAUDE.md` — здесь только про проект.

> Параллельно существует `AGENTS.md` (рус) и `GEMINI.md` (англ) с теми же правилами.
> При изменении архитектуры/флоу/контрактов синхронно обновляй и `AGENTS.md` — это жёсткое правило проекта.

> **Этот сервис — часть системы BSKT.IMP** (роль — integrator). Соседи: API на Laravel и dashboard на Vue.
> Межсервисные контракты и карта живут в Obsidian-vault **IMP** (через MCP `obsidian`, если подключён): `system-map`, `contracts/db-schema`, `contracts/aggregator-webhook`, `adr/`.
> **Перед изменением того, что пересекает границу сервиса, — сверься с vault:**
> - меняешь **схему БД** (миграции): API читает ту же PostgreSQL напрямую, смена формы таблиц может молча сломать его → add-only / согласование (см. `contracts/db-schema`, `adr/0001`);
> - меняешь **webhook `game-imported`** (payload/endpoint): его потребитель — Laravel-API (см. `contracts/aggregator-webhook`).
> Внутренние изменения, не пересекающие границу сервиса, vault не трогают.

---

## 1. Что это

Сервис собирает баскетбольную статистику из внешних провайдеров и сохраняет в PostgreSQL:
игры, командную статистику, индивидуальную статистику игроков и служебные watermark-метки опроса.

Источники данных:
- `API_NBA` — полная поддержка (статистика + био игроков);
- `INFOBASKET` — игры и статистика;
- `SPORTOTEKA` — игры и статистика;
- `CDN_NBA` — заглушка в коде, не активен.

## 2. Стек

- Go **1.25**, гексагональная архитектура (ports/adapters)
- PostgreSQL + GORM (ORM), goose (миграции)
- `golang.org/x/time/rate` — rate limiting для внешних API
- gomock + testify (assert/suite), линт — golangci-lint (с testifylint)
- Логгер — внешний `github.com/Consolushka/golang.composite_logger` (консоль / файл / Telegram)
- Docker + docker-compose для PostgreSQL и контейнеризации приложения

## 3. Архитектура (слои)

Поток зависимостей строго внутрь: `service` зависит от `ports` (интерфейсов), а не от конкретных адаптеров.

- `app/internal/core/<domain>` — доменные модели (games, players, teams, leagues, tournaments, tournament_poll_logs). Модели разных доменов не смешивать в одном пакете.
- `app/internal/ports` — интерфейсы-контракты (репозитории, `StatsProvider`, aggregator, cached_remote_resource).
- `app/internal/adapters` — реализации: GORM-репозитории (`<name>_repo/gorm.go`), стат-провайдеры, кеш.
- `app/internal/service` — use-case и оркестрация (orchestrator, processor, persistence, scheduler, providers).
- `app/internal/infra` — низкоуровневые HTTP-клиенты внешних API + "чистые" трансформеры данных + фабрики (config, logger).
- `app/database/migrations` — SQL-миграции goose.
- `app/pkg` — переиспользуемые утилиты (`dbtest`, `statsutil`, `http`).

Точки входа (`app/cmd`):
- `scheduler/main.go` — основной фоновый процесс (распределённый опрос турниров);
- `debug-server/main.go` — HTTP-сервер на `:8080` для ручного триггера обработки;
- `create-tournament/main.go` — CLI-команда регистрации турнира (лига + привязка к провайдеру + даты).

## 4. Флоу обработки данных

1. **Scheduler** периодически (`SCHEDULER_REFRESH_INTERVAL_MINUTES`, по умолч. 5) обновляет список активных турниров: поднимает горутину-воркер на каждый новый турнир и гасит воркеры деактивированных — без рестарта сервиса.
2. **Distributed workers**: на каждый турнир — своя горутина. Старты разнесены во времени на `SCHEDULER_STAGGER_INTERVAL_MINUTES` ("шахматный старт"), чтобы размазать нагрузку на API и БД. Каждый воркер опрашивает свой турнир раз в `SCHEDULER_POLL_INTERVAL` минут.
3. Интервал следующего опроса берётся из `interval_end` последней успешной записи в `tournament_poll_logs` (watermark по `tournament_id`).
4. **Orchestrator / processor** тянут и обогащают данные через `StatsProvider`. Контракт намеренно разделён: лёгкий листинг (`GetGamesStatsByPeriod`) и тяжёлое обогащение (`EnrichGameStats`) — ради лимитов API.
5. **Обработка игроков в две фазы**:
   - *Discovery* — есть ли игрок в локальной БД;
   - *Ingestion* — если нет/данные неполные, запросить био у провайдера (`GetPlayerBio`).
6. **Persistence** сохраняет данные; результат опроса (статус, время, число игр, ошибки) пишется в `tournament_poll_logs`.

Инварианты данных:
- Трансформеры в `infra` — **чистые функции**: только маппинг, никаких сетевых вызовов.
- Все процентные показатели (FG%, и т.п.) нормализуются к диапазону **`0.0–1.0`** на уровне трансформера, независимо от формата провайдера. На уровне БД диапазон гарантируется `CHECK`-констрейнтами (см. `game_team_player_stats`).
- Rate limiting встроен в infra-клиенты, настраивается через `*_RATE_LIMIT_PER_MINUTE`.
- `ListActive` не возвращает завершённые турниры (`end_at < now`), чтобы не тратить лимиты API. `end_at IS NULL` означает «даты неизвестны» — такой турнир считается **активным**.

### Регистрация турнира

Use-case `service.TournamentRegistrar.Create` (CLI `create-tournament`) создаёт турнир: гарантирует лигу (`LeaguesRepo.FirstOrCreate`), определяет даты сезона и атомарно пишет турнир + provider-маппинг (`TournamentsRepo.Create`). Даты по приоритету: **явные (флаги `-start/-end`) → провайдер → предупреждение + `NULL`**. «Провайдер умеет отдавать даты» — это опциональный порт `TournamentPeriodProvider` (его реализует только `API_BASKETBALL` через `/leagues`); проверяется приведением типа, не реализован → warning. `ProviderFactory` инжектируется в use-case (в проде — `providers.NewProvider`), чтобы он был тестируем без HTTP.

## 5. Запуск и команды (Makefile)

```bash
make setup          # cp .example.env .env
make db-start       # PostgreSQL в фоне (docker), db-stop/db-down — остановка
make migrate        # goose up
make run-scheduler  # основной процесс
make run-debug      # debug-сервер на :8080
make test           # go test ./...
make test-with-coverage  # с покрытием, моки исключаются из отчёта
make lint           # golangci-lint, lint-fix — с автофиксом
make create-migration name=add_some_column
```

Ключевые env (шаблон — `.example.env`): `DB_*`, `GOOSE_*`, `API_SPORT_API_KEY` (для API_NBA),
`SCHEDULER_POLL_INTERVAL` / `SCHEDULER_STAGGER_INTERVAL_MINUTES` / `SCHEDULER_REFRESH_INTERVAL_MINUTES`,
`*_RATE_LIMIT_PER_MINUTE`, `LOGGER_*`.

Debug API:
- `GET /health`
- `GET /process/all?from=YYYY-MM-DD&to=YYYY-MM-DD`
- `GET /process/tournament?id=N&from=...&to=...`
- без `from`/`to` — период с начала текущих UTC-суток.

Регистрация турнира (CLI):
```bash
go run ./app/cmd/create-tournament \
  -name "Euroleague 2025-2026" -league-name Euroleague -league-alias euroleague \
  -provider API_BASKETBALL -external-id 120 -season 2025
```
Обязательны `-name`, `-league-name`, `-league-alias`, `-provider`. Провайдерские params — JSON-флагом `-params` (напр. `'{"leadHost":"reg","year":2024}'`). Явные `-start/-end` (YYYY-MM-DD) перекрывают даты от провайдера.

Деплой: `docker/dockhost/Dockerfile`. `startup.sh` прогоняет `go test ./...` перед стартом бинарника — **падающие тесты блокируют запуск**.

## 6. БД и миграции

- Миграции — `app/database/migrations` (goose, `make create-migration name=...`).
- Перед изменением моделей проверяй существующие индексы/уникальности.
- Критерий "игра уже существует" в коде должен соответствовать реальному unique-индексу в БД.
- Для процентных колонок — `CHECK`-констрейнты на диапазон `[0, 1]`.
- Миграции обязаны накатываться на чистую базу.

## 7. Тесты

- **Репозитории (GORM)** — SQLite in-memory через `app/pkg/dbtest`. Паттерн: `testify/suite`, БД и миграции поднимаются один раз в `SetupSuite`; каждый тест — внутри транзакции (`db.Begin()` в `SetupTest`, `Rollback()` в `TearDownTest`); репозиторий инициализируется транзакционным объектом (`s.tx`).
- **Провайдеры** — unit-тесты с `MockClientInterface`: успешный маппинг, фильтрация игр по статусу, ошибки API на этапах листинга/деталей, работа с кешем.
- Для проверки ошибок — **`s.Require().NoError(err)`**, не `s.NoError` (правило `testifylint: require-error`).
- Любая новая ветвящаяся логика покрывается: happy-path, error-path, continue/skip-ветки, взаимодействие с моками.
- Меняешь сигнатуру интерфейса в `ports` → обнови все моки и тесты, добейся компиляции и зелёного `make lint`.

## 8. Правила кода

- Не нарушать слои; не тянуть `adapters`/`infra` напрямую в `core`/`service`.
- Новый флоу/зависимость — сначала порт (интерфейс), потом адаптер.
- Репозитории: интерфейс в `ports/*_repo.go`, реализация в `adapters/<name>_repo/gorm.go`, структура `Gorm`, ресивер `g`, конструктор `NewGormRepo(...)`.
- Именование методов репозиториев без суффикса сущности: `FirstOrCreate`, `Exists`, `Get`, `ListActive`, `ListBy...`.
- **Compile-time проверка контракта**: для каждой конкретной реализации порта добавляй assertion, чтобы расхождение с интерфейсом ломало компиляцию, а не всплывало в рантайме:

  ```go
  // в adapters/games_repo/gorm.go
  var _ ports.GamesRepo = (*Gorm)(nil)
  ```

  Форма `(*Type)(nil)` работает и для value-, и для pointer-ресиверов. То же — для адаптеров провайдеров (`var _ ports.StatsProvider = (*Adapter)(nil)`).
- Долгоживущие горутины (воркеры scheduler) защищать `defer composite_logger.Recover(ctx)`.
- Избегать гонок в замыканиях горутин.
- Логи информативные, не шумные: не логировать тяжёлые структуры целиком; в error-логах — только поля, по которым ищется/проверяется сущность; не дублировать логи между слоями persistence и service.

## 9. Definition of Done

Задача считается готовой, только когда:
1. `make lint` — зелёный (golangci-lint + testifylint).
2. `make test` — зелёный (это же гейт в `startup.sh` перед деплоем).
3. Интерфейсы и их реализации/моки синхронны (поменял `ports` — поменял моки).
4. Новая ветвящаяся логика покрыта тестами.
5. Миграции накатываются на чистую базу.
6. При изменении архитектуры/флоу/контрактов/БД/правил — обновлён `AGENTS.md` (и `README.md`, если менялось поведение).

## 10. Git / релизный флоу

- Фичи — ветки `feature/<kebab-описание>`.
- Релиз — ветка `release/X.Y`, вливается в `main` через PR с заголовком `Release X.Y (#NN)`.
- Срочные правки — коммиты с префиксом `hotfix:` / `hotfix(scope):`.
- Стиль сообщений — conventional-ish: `feat:`, `feat(scope):`, `refactor:`, `hotfix:`.
- Коммитить/пушить — только по явной просьбе. Если на `main` — сначала ветка.

## 11. Карта "куда смотреть"

| Задача | Где трогать |
|---|---|
| Новый стат-провайдер | `infra/<provider>` (клиент + чистый трансформер) → `adapters/stats_provider` → интерфейс в `ports` → регистрация в `service/providers` |
| Новое поле/таблица | `database/migrations` (+ индексы/`CHECK`) → модель в `core/<domain>` → репозиторий в `adapters/<name>_repo` |
| Изменение логики опроса | `service/scheduler` + `service` (orchestrator/processor) |
| Лимиты внешних API | rate limiter в `infra/<provider>` + `*_RATE_LIMIT_PER_MINUTE` в `.env` |
| Регистрация турнира | `service/tournament_registrar.go` + CLI `app/cmd/create-tournament` + репозитории `leagues_repo`/`tournaments_repo` |
| Даты турнира от провайдера | порт `ports.TournamentPeriodProvider` + реализация в `adapters/stats_provider/<provider>.go` |
| Ручной прогон | `app/cmd/debug-server` |
| Конфиг/env | `infra/config`, `.example.env` |
| Логирование | `infra/logger` (фабрика над composite_logger) |
