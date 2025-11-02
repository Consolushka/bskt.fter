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

type ProcessNotUrgentEuropeanTournaments struct {
	task func(from, to time.Time) error
}

func (p ProcessNotUrgentEuropeanTournaments) GetPeriodicity() time.Duration {
	return 24 * time.Hour
}

func (p ProcessNotUrgentEuropeanTournaments) GetName() string {
	return "ProcessByPeriod not urgent european tournaments task"
}

func (p ProcessNotUrgentEuropeanTournaments) Execute(lastExecutedAt time.Time) error {
	return p.task(lastExecutedAt, time.Now())
}

func NewProcessNotUrgentEuropeanTournaments(db *gorm.DB) ports.ExecutableByScheduler {
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
		players_repo.NewGormRepo(db),
	)
	return ProcessNotUrgentEuropeanTournaments{
		task: orchestrator.ProcessNotUrgentEuropeanTournaments,
	}
}
