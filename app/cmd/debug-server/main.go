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

	// Инициализируем оркестратор, как в твоем старом локальном main.go
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(
			games_repo.NewGormRepo(db),
			teams_repo.NewGormRepo(db),
			players_repo.NewGormRepo(db),
		),
		tournaments_repo.NewGormRepo(db),
		players_repo.NewGormRepo(db),
	)

	// Вспомогательная функция для обработки дат
	parseDates := func(r *http.Request) (time.Time, time.Time, error) {
		layout := "2006-01-02"
		fromStr := r.URL.Query().Get("from")
		toStr := r.URL.Query().Get("to")

		if fromStr == "" || toStr == "" {
			return time.Time{}, time.Time{}, fmt.Errorf("missing 'from' or 'to' parameters (format: YYYY-MM-DD)")
		}

		from, err := time.Parse(layout, fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'from' date")
		}

		to, err := time.Parse(layout, toStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'to' date")
		}

		// Устанавливаем время для захвата всего дня, аналогично твоему коду
		return from.UTC(), to.UTC().Add(24*time.Hour - time.Nanosecond), nil
	}

	// Эндпоинты
	http.HandleFunc("/process/american", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDates(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logger.Info(fmt.Sprintf("Manual trigger: American Tournaments from %s to %s", from, to), nil)
		err = orchestrator.ProcessAmericanTournaments(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK: American tournaments processed"))
	})

	http.HandleFunc("/process/european-urgent", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDates(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = orchestrator.ProcessUrgentEuropeanTournaments(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK: Urgent European tournaments processed"))
	})

	http.HandleFunc("/process/european-not-urgent", func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDates(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = orchestrator.ProcessNotUrgentEuropeanTournaments(from, to)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("OK: Not urgent European tournaments processed"))
	})

	logger.Info("Debug API server started on :8080", nil)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
