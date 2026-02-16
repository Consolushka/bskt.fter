# Gemini CLI Context: IMP (Basketball Statistics Collector)

This project, named **IMP**, is a Go-based service designed to collect basketball match statistics from various external providers and persist them in a PostgreSQL database.

## Project Overview

*   **Purpose:** Collects games, team stats, and player stats from providers like `API_NBA`, `INFOBASKET`, and `SPORTOTEKA`.
*   **Architecture:** Follows **Hexagonal Architecture (Ports and Adapters)**:
    *   `app/internal/core`: Domain entities and models.
    *   `app/internal/ports`: Interfaces (contracts) for repositories and providers.
    *   `app/internal/adapters`: Implementations of ports (GORM for DB, specific provider logic).
    *   `app/internal/service`: Business logic, orchestration (scheduler, processor, persistence).
    *   `app/internal/infra`: Low-level HTTP clients and transformers for external APIs.
*   **Main Components:**
    *   **Scheduler:** A background worker (`app/cmd/scheduler`) that manages distributed workers for each tournament with staggered start times.
    *   **Debug Server:** A small HTTP server (`app/cmd/debug-server`) for manual triggers and debugging.
    *   **Persistence:** Uses GORM to interact with PostgreSQL.
    ## Technologies

*   **Language:** Go 1.23.4
*   **Database:** PostgreSQL
*   **ORM:** GORM
*   **Logging:** Logrus (with support for console, file, and Telegram)
*   **Testing:** Testify, GoMock
*   **Rate Limiting:** golang.org/x/time/rate
*   **Migrations:** Goose
*   **Linting:** golangci-lint

## Building and Running

The project uses a `Makefile` for common tasks:

*   **Setup:** `make setup` (copies `.example.env` to `.env`).
*   **Infrastructure:**
    *   `make start`: Run PostgreSQL in background (Docker).
    *   `make stop`: Stop PostgreSQL.
*   **Running:**
    *   `make run-scheduler`: Start the background worker.
    *   `make run-debug`: Start the debug HTTP server (default port `:8080`).
*   **Database:**
    *   `make migrate`: Apply database migrations.
    *   `make create-migration name=<name>`: Create a new migration file.
*   **Quality:**
    *   `make test-with-coverage`: Run tests and generate coverage report.
    *   `make lint`: Run `golangci-lint`.
    *   `make lint-fix`: Run `golangci-lint` with auto-fix.

## Development Conventions

*   **Strict Layering:** Never import `adapters` or `infra` directly into `core` or `service`. Use `ports` (interfaces).
*   **Naming:**
    *   Repositories should have interfaces in `ports` and implementations in `adapters/<name>_repo/gorm.go`.
    *   Constructors should follow the `NewGormRepo(...)` pattern.
    *   Methods in single-entity repositories should avoid entity suffixes (e.g., `FirstOrCreate` instead of `FirstOrCreateTeam`).
    *   Receiver names in GORM adapters should be `g`.
*   **Mocking:** Use `gomock` for testing. Mocks should be updated whenever an interface in `ports` changes.
*   **Error Handling:** Use standard library `errors`. Avoid `github.com/pkg/errors`.
*   **Linting:** Adhere to `.golangci.yml` rules. Shadowing is checked, and cyclomatic complexity is limited to 20.
*   **Documentation:** Always update `AGENTS.md` when architectural or significant logic changes are made.
*   **Environment:** Use `.env` for configuration. Mandatory variables include `DB_*` settings, provider API keys (e.g., `API_SPORT_API_KEY`), rate limits (`*_RATE_LIMIT_PER_MINUTE`), and scheduler settings (`SCHEDULER_POLL_INTERVAL`, `SCHEDULER_STAGGER_INTERVAL_MINUTES`).

## Directory Structure Highlights

*   `app/cmd/`: Entry points for the application.
*   `app/database/migrations/`: SQL migration files managed by Goose.
*   `app/internal/adapters/`: Database and external provider implementations.
*   `app/internal/core/`: Domain models (e.g., `games`, `players`, `teams`).
*   `app/internal/ports/`: Interface definitions.
*   `app/internal/service/`: Core business logic and orchestrators.
