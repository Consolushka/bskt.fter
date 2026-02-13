package executable_by_scheduler

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"time"

	"gorm.io/gorm"
)

type ProcessUrgentEuropeanTournaments struct {
	task func(from, to time.Time) error
}

func (p ProcessUrgentEuropeanTournaments) GetPeriodicity() time.Duration {
	return 8 * time.Hour
}

func (p ProcessUrgentEuropeanTournaments) GetName() string {
	return "ProcessByPeriod urgent european tournaments task"
}

func (p ProcessUrgentEuropeanTournaments) Execute(lastExecutedAt time.Time) error {
	return p.task(lastExecutedAt, time.Now())
}

func NewProcessUrgentEuropeanTournaments(db *gorm.DB) ports.ExecutableByScheduler {
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
		players_repo.NewGormRepo(db),
		games_repo.NewGormRepo(db),
	)
	return ProcessUrgentEuropeanTournaments{
		task: orchestrator.ProcessUrgentEuropeanTournaments,
	}
}
