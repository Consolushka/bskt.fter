# CLAUDE.md

Технический контекст проекта **IMP** для работы в этом репозитории.
Роль, стиль общения и инженерные принципы заданы глобально в `~/.claude/CLAUDE.md` — здесь только про проект.

> Параллельно существует `AGENTS.md` (рус) и `GEMINI.md` (англ) с теми же правилами.
> При изменении архитектуры/флоу/контрактов синхронно обновляй и `AGENTS.md` — это жёсткое правило проекта.

> **Argus** — этот сервис (роль integrator в системе BSKT.IMP). Соседи: API на Laravel и dashboard на Vue.
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
- `ESPN_NBA` — полная поддержка (игры, статистика, био игроков); основной источник для NBA, публичный API без ключа;
- `API_NBA` — полная поддержка (статистика + био игроков);
- `INFOBASKET` — игры и статистика;
- `SPORTOTEKA` — игры и статистика.

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
5. **Команды и название игры**: до проверки существования игры процессор находит обе команды в базе (по id провайдера, запасной путь — по имени) и собирает `games.title` из их алиасов.
6. **Обработка игроков в две фазы**:
   - *Discovery* — есть ли игрок в локальной БД;
   - *Ingestion* — если нет/данные неполные, запросить био у провайдера (`GetPlayerBio`).
7. **Persistence** сохраняет данные; результат опроса (статус, время, число игр, ошибки) пишется в `tournament_poll_logs`.

Инварианты данных:
- Трансформеры в `infra` — **чистые функции**: только маппинг, никаких сетевых вызовов.
- Все процентные показатели (FG%, и т.п.) нормализуются к диапазону **`0.0–1.0`** на уровне трансформера, независимо от формата провайдера. На уровне БД диапазон гарантируется `CHECK`-констрейнтами (см. `game_team_player_stats`).
- `game_team_stats.is_home` выставляет persistence по роли (`HomeTeamStat`/`AwayTeamStat`), не трансформер. Без значения по умолчанию; частичный unique-индекс — не больше одного хозяина в игре. Это контракт с API (`contracts/db-schema`).
- Rate limiting встроен в infra-клиенты, настраивается через `*_RATE_LIMIT_PER_MINUTE`.
- `ListActive` не возвращает турниры, закончившиеся больше 14 дней назад (`end_at < now − 14 дней`), чтобы не тратить лимиты API. Первые 14 дней после `end_at` опрос продолжается: игры плей-офф и переносы сохраняются и дают warning, но в лог опроса пишется `next_poll_at = NULL` — завершённому турниру следующее обновление не обещаем. `end_at IS NULL` означает «даты неизвестны» — такой турнир считается **активным**.
- Даты турнира Argus сам не меняет. Сохранённая игра раньше `start_at` или позже `end_at` (с запасом 36 часов на часовые пояса) даёт warning, один на опрос турнира; решение продлевать даты — за владельцем. Воркер перечитывает турнир перед каждым опросом, поэтому ручная правка дат подхватывается без рестарта.
- Коды команд NBA канонизируются к кодам api-sports (`GS→GSW`, `NO→NOP`, `NY→NYK`, `SA→SAS`, `UTAH→UTA`, `WSH→WAS`) в трансформере ESPN: алиас берётся из кода при создании команды, `teams.alias` задан первым провайдером, и ESPN не должен заводить команду с другим кодом.
- `teams.alias` пишется только при создании команды и дальше не перезаписывается (ручные правки не откатываются). Сокращение INFOBASKET/SPORTOTEKA из 6 символов, совпадающее с началом имени, — это обрезанное имя: `teams.ProviderAlias` его отбрасывает, и алиас генерирует `TeamModel.AutoGenerateAlias` (без кавычек и префикса `BK`/`BC`/`БК`; имя до 6 символов целиком, иначе 3 буквы).
- ESPN: трансформер сверяет сумму очков и минут по игрокам с итогом команды; расхождение — ошибка `EnrichGameStats`, игра не сохраняется (лучше отсутствие игры, чем неверная статистика).
- Дни ESPN — календарные дни US Eastern: адаптер переводит период опроса в `America/New_York` и запрашивает scoreboard по одному дню.
- **Название игры**: `games.title` = `<teams.alias хозяев> - <teams.alias гостей>`. Собирает процессор (`GameStatEntity.AssignTeams`) из алиасов в базе до проверки существования игры; трансформеры его не заполняют.
- **Сопоставление команд**: команда ищется по `(scope, external_id)` в `team_external_ids`, запасной путь — `(name, home_town)` (если провайдер не отдал город — только по имени), после него связь записывается. `scope` — имя провайдера, у Infobasket — `INFOBASKET:<leadHost>`. Команды создаёт только процессор. Перед регистрацией турнира новой лиги Sportoteka — сверка запросом ниже: у каждой возвращённой строки имя и город должны соответствовать команде с этим id в календаре новой лиги; при расхождении турнир не регистрировать, пока адаптер Sportoteka не получит собственный `scope`.

  ```sql
  SELECT e.external_id, t.id AS team_id, t.name, t.home_town, t.alias
  FROM team_external_ids e
  JOIN teams t ON t.id = e.team_id
  WHERE e.scope = 'SPORTOTEKA' AND e.external_id IN (/* teamId из календаря новой лиги */);
  ```
- **Новая команда**: о каждой новой строке `teams` пишется warning `New team created` — владелец проверяет алиас.
- **Ручная правка `teams.alias`**: только одной транзакцией с `UPDATE games.title` игр этой команды: название входит в ключ игры, и без пересчёта игры команды из окна повторного опроса запишутся второй раз. Править между опросами турниров этой команды и после правки проверить, что игра не записалась второй раз со старым названием: воркер мог прочитать старый алиас до транзакции.

  ```sql
  BEGIN;
  UPDATE teams SET alias = :new_alias WHERE id = :team_id;
  UPDATE games g
  SET title = h.alias || ' - ' || a.alias
  FROM game_team_stats hs
  JOIN teams h             ON h.id = hs.team_id
  JOIN game_team_stats aws ON aws.game_id = hs.game_id AND NOT aws.is_home
  JOIN teams a             ON a.id = aws.team_id
  WHERE hs.game_id = g.id AND hs.is_home
    AND :team_id IN (hs.team_id, aws.team_id);
  COMMIT;
  ```
- **Ошибочная связь команды с id провайдера**: сбрасывается запросом ниже; следующий опрос свяжет id заново. Команды, на которые в одном `scope` ведут два id (слияние двух команд провайдера в одну), ищутся вторым запросом.

  ```sql
  DELETE FROM team_external_ids WHERE scope = :scope AND external_id = :external_id;
  ```

  ```sql
  SELECT e.scope, e.team_id, t.name, t.home_town, array_agg(e.external_id ORDER BY e.external_id) AS external_ids
  FROM team_external_ids e
  JOIN teams t ON t.id = e.team_id
  GROUP BY e.scope, e.team_id, t.name, t.home_town
  HAVING count(*) > 1;
  ```

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

Ключевые env (шаблон — `.example.env`): `DB_*`, `GOOSE_*`, `API_SPORT_API_KEY` (для API_NBA и API_BASKETBALL),
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
| Новый стат-провайдер | `infra/<provider>` (клиент + чистый трансформер) → `adapters/stats_provider` → интерфейс в `ports` → регистрация в `service/providers` (образец пути «листинг → протокол» и тестов на фикстурах — `infra/espn`) |
| Новое поле/таблица | `database/migrations` (+ индексы/`CHECK`) → модель в `core/<domain>` → репозиторий в `adapters/<name>_repo` |
| Изменение логики опроса | `service/scheduler` + `service` (orchestrator/processor) |
| Лимиты внешних API | rate limiter в `infra/<provider>` + `*_RATE_LIMIT_PER_MINUTE` в `.env` |
| Регистрация турнира | `service/tournament_registrar.go` + CLI `app/cmd/create-tournament` + репозитории `leagues_repo`/`tournaments_repo` |
| Даты турнира от провайдера | порт `ports.TournamentPeriodProvider` + реализация в `adapters/stats_provider/<provider>.go` |
| Ручной прогон | `app/cmd/debug-server` |
| Конфиг/env | `infra/config`, `.example.env` |
| Сопоставление команд и название игры | `service/tournament_processor.go` (`resolveTeam`) + `core/games/entities.go` (`AssignTeams`) + `adapters/teams_repo` + порт `ports.TeamIdScopeProvider` |
| Логирование | `infra/logger` (фабрика над composite_logger) |
