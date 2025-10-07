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

type ProcessAmericanTournaments struct {
	task func() error
}

func (p ProcessAmericanTournaments) GetPeriodicity() time.Duration {
	return 24 * time.Hour
}

func (p ProcessAmericanTournaments) GetName() string {
	return "Process american tournaments task"
}

func (p ProcessAmericanTournaments) ShouldBeExecutedAt() time.Time {
	now := time.Now()
	loc := time.FixedZone("UTC-7", -7*60*60)

	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, now.Nanosecond(), loc)
}

func (p ProcessAmericanTournaments) Execute() error {
	return p.task()
}

func NewProcessAmericanTournaments(db *gorm.DB) ports.ExecutableByScheduler {
	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
	)
	return ProcessAmericanTournaments{
		task: orchestrator.ProcessAmericanTournaments,
	}
}
