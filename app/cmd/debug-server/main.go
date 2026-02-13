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
	"time"

	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC

	godotenv.Load()
	logger.Init(logger.BuildLoggers())

	database.OpenDbConnection()
	db := database.GetDB()

	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(
			games_repo.NewGormRepo(db),
			teams_repo.NewGormRepo(db),
			players_repo.NewGormRepo(db),
		),
		tournaments_repo.NewGormRepo(db),
		players_repo.NewGormRepo(db),
		games_repo.NewGormRepo(db),
	)

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/process/american", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parsePeriod(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := orchestrator.ProcessAmericanTournaments(from, to); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("ok american from=%s to=%s", from.Format(time.RFC3339), to.Format(time.RFC3339))))
	})

	http.HandleFunc("/process/european-urgent", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parsePeriod(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := orchestrator.ProcessUrgentEuropeanTournaments(from, to); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("ok european-urgent from=%s to=%s", from.Format(time.RFC3339), to.Format(time.RFC3339))))
	})

	http.HandleFunc("/process/european-not-urgent", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parsePeriod(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := orchestrator.ProcessNotUrgentEuropeanTournaments(from, to); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("ok european-not-urgent from=%s to=%s", from.Format(time.RFC3339), to.Format(time.RFC3339))))
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
