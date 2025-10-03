package executable_by_scheduler

import (
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"time"
)

type ProcessAmericanTournaments struct {
	task func() error
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

func NewProcessNightlyTournaments(orchestrator *service.TournamentsOrchestrator) ports.ExecutableByScheduler {
	return ProcessAmericanTournaments{
		task: orchestrator.ProcessAmericanTournaments,
	}
}
