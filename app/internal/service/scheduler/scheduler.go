package scheduler

import (
	"IMP/app/internal/adapters/executable_by_scheduler"
	"IMP/app/internal/adapters/scheduler_repo"
	"IMP/app/internal/core/scheduler"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"context"
	"errors"
	"sync"
	"time"

	"gorm.io/gorm"
)

func Handle(db *gorm.DB) {
	schedulerRepo := scheduler_repo.NewGormRepo(db)

	tasks, err := schedulerRepo.TasksList()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		go func() {
			err = handleTask(task, db, ctx)
			if err != nil {
				logger.Error("Error executing task", map[string]interface{}{
					"task":  task,
					"error": err,
				})
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func handleTask(taskModel scheduler.ScheduledTaskModel, db *gorm.DB, ctx context.Context) error {
	schedulerRepo := scheduler_repo.NewGormRepo(db)

	for {
		task := matchExecutableAdapterByTaskType(taskModel.Type, db)
		if task == nil {
			return errors.New("Couldn't find executable adapter for task type: " + taskModel.Type)
		}

		now := time.Now()

		var sleepDuration time.Duration

		if now.Before(taskModel.NextExecutionAt) {
			sleepDuration = taskModel.NextExecutionAt.Sub(now)

			logger.Info(task.GetName()+" will be executed at "+taskModel.NextExecutionAt.Format("02-01-2006 15:04"), map[string]interface{}{
				"taskModel": taskModel,
			})
		} else {
			logger.Info(task.GetName()+" should been executed at "+taskModel.NextExecutionAt.Format("02-01-2006 15:04")+". Executing...", map[string]interface{}{})

			sleepDuration = taskModel.NextExecutionAt.Sub(now)
		}

		timer := time.NewTimer(sleepDuration)

		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			err := task.Execute(taskModel.LastExecutedAt)

			if err != nil {
				logger.Error("Error while processing tournament games", map[string]interface{}{
					"error": err,
				})
			}

			taskModel, err = schedulerRepo.RescheduleTask(taskModel.Id, taskModel.NextExecutionAt.Add(task.GetPeriodicity()))
			if err != nil {
				logger.Warn("Couldn't reschedule taskModel", map[string]interface{}{
					"error": err,
				})
			}
		}
	}
}

func matchExecutableAdapterByTaskType(taskType string, db *gorm.DB) ports.ExecutableByScheduler {
	switch taskType {
	case ProcessAmericanTournamentsTaskType:
		return executable_by_scheduler.NewProcessAmericanTournaments(db)
	case ProcessNotUrgentEuropeanTournamentsTaskType:
		return executable_by_scheduler.NewProcessNotUrgentEuropeanTournaments(db)
	case ProcessUrgentEuropeanTournamentsTaskType:
		return executable_by_scheduler.NewProcessUrgentEuropeanTournaments(db)
	default:
		return nil
	}
}
