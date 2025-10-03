package executable_by_scheduler

import (
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/pkg/time_utils"
	"time"
)

type ProcessTodayTournaments struct {
	task func() error
}

func (p ProcessTodayTournaments) GetName() string {
	return "Process urgent european tournaments task"
}

func (p ProcessTodayTournaments) ShouldBeExecutedAt() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, now.Nanosecond(), time_utils.MoscowTZ)
}

func (p ProcessTodayTournaments) Execute() error {
	return p.task()
}

func NewProcessTodayTournaments(orchestrator *service.TournamentsOrchestrator) ports.ExecutableByScheduler {
	return ProcessTodayTournaments{
		task: orchestrator.ProcessUrgentEuropeanTournaments,
	}
}
