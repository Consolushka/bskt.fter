package executable_by_scheduler

import (
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/pkg/time_utils"
	"time"
)

type ProcessYesterdayDailyTournaments struct {
	task func() error
}

func (p ProcessYesterdayDailyTournaments) GetName() string {
	return "Process not urgent european tournaments task"
}

func (p ProcessYesterdayDailyTournaments) ShouldBeExecutedAt() time.Time {
	now := time.Now()

	return time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, now.Nanosecond(), time_utils.MoscowTZ)
}

func (p ProcessYesterdayDailyTournaments) Execute() error {
	return p.task()
}

func NewProcessYesterdayDailyTournaments(orchestrator *service.TournamentsOrchestrator) ports.ExecutableByScheduler {
	return ProcessYesterdayDailyTournaments{
		task: orchestrator.ProcessNotUrgentEuropeanTournaments,
	}
}
