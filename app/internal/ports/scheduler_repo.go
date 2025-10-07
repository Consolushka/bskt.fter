package ports

import (
	"IMP/app/internal/core/scheduler"
	"time"
)

type SchedulerRepo interface {
	TasksList() ([]scheduler.ScheduledTaskModel, error)
	RescheduleTask(id uint, nextExecutionAt time.Time) (scheduler.ScheduledTaskModel, error)
}
