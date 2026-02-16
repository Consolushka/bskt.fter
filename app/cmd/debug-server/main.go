package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/service"
	"IMP/app/pkg/logger"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// 2.2
// todo: Ввести единый provider rate limiter слой
// 2.3
// todo: gap between every poll. for example betwen every single tournament delay is 30, but betwen different is like 5 minutes.
// 2.4
// todo: test before deploy
// 2.5
// todo: db integration tests
// 2.6
// todo: мб логгер в defer?
// todo: Разделить pipeline на discovery и ingestion явно
// 3.0
// todo: Добавить provider health/state в БД
// todo: Сделать observability минимально production-ready
func main() {
	time.Local = time.UTC

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	logger.Init(logger.BuildLoggers())

	database.OpenDbConnection()
	db := database.GetDB()

	tr := tournaments_repo.NewGormRepo(db)
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(
			games_repo.NewGormRepo(db),
			teams_repo.NewGormRepo(db),
			players_repo.NewGormRepo(db),
		),
		tr,
		players_repo.NewGormRepo(db),
		games_repo.NewGormRepo(db),
	)

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			logger.Warn("Failed to write health check response", map[string]interface{}{"error": err})
			return
		}
	})

	http.HandleFunc("/process/all", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parsePeriod(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = orchestrator.ProcessAll(from, to); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("ok all from=%s to=%s", from.Format(time.RFC3339), to.Format(time.RFC3339))))
		if err != nil {
			logger.Warn("Failed to write all response", map[string]interface{}{"error": err})
			return
		}
	})

	http.HandleFunc("/process/tournament", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parsePeriod(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			http.Error(w, "invalid tournament id", http.StatusBadRequest)
			return
		}

		tournament, err := tr.Get(uint(id))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get tournament: %v", err), http.StatusNotFound)
			return
		}

		if err = orchestrator.ProcessTournament(tournament, from, to); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("ok tournament=%d from=%s to=%s", id, from.Format(time.RFC3339), to.Format(time.RFC3339))))
		if err != nil {
			logger.Warn("Failed to write tournament response", map[string]interface{}{"error": err})
		}
	})

	logger.Info("Debug server started", map[string]interface{}{"port": 8080})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func parsePeriod(r *http.Request) (time.Time, time.Time, error) {
	const layout = "2006-01-02"
	fromRaw := r.URL.Query().Get("from")
	toRaw := r.URL.Query().Get("to")

	now := time.Now().UTC()
	defaultFrom := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	defaultTo := now

	if fromRaw == "" && toRaw == "" {
		return defaultFrom, defaultTo, nil
	}

	if fromRaw == "" || toRaw == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("both from and to are required in format YYYY-MM-DD")
	}

	from, err := time.Parse(layout, fromRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid from date")
	}

	to, err := time.Parse(layout, toRaw)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid to date")
	}

	from = from.UTC()
	to = to.UTC().Add(24*time.Hour - time.Nanosecond)
	if to.Before(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("to must be greater or equal to from")
	}

	return from, to, nil
}
