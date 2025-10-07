package executable_by_scheduler

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/pkg/time_utils"
	"time"

	"gorm.io/gorm"
)

type ProcessNotUrgentEuropeanTournaments struct {
	task func() error
}

func (p ProcessNotUrgentEuropeanTournaments) GetPeriodicity() time.Duration {
	return 24 * time.Hour
}

func (p ProcessNotUrgentEuropeanTournaments) GetName() string {
	return "Process not urgent european tournaments task"
}

func (p ProcessNotUrgentEuropeanTournaments) ShouldBeExecutedAt() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, now.Nanosecond(), time_utils.MoscowTZ)
}

func (p ProcessNotUrgentEuropeanTournaments) Execute() error {
	return p.task()
}

func NewProcessNotUrgentEuropeanTournaments(db *gorm.DB) ports.ExecutableByScheduler {
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
	)
	return ProcessNotUrgentEuropeanTournaments{
		task: orchestrator.ProcessNotUrgentEuropeanTournaments,
	}
}
