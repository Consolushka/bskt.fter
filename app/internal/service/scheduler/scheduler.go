package scheduler

import (
	"IMP/app/internal/adapters/executable_by_scheduler"
	"IMP/app/internal/adapters/poll_watermarks_repo"
	"IMP/app/internal/core/poll_watermarks"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"os"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

func Handle(db *gorm.DB) {
	pollIntervalString := os.Getenv("SCHEDULER_POLL_INTERVAL")
	pollIntervalInMinutes, err := strconv.Atoi(pollIntervalString)
	if err != nil || pollIntervalInMinutes <= 0 {
		logger.Warn("Couldn't load SCHEDULER_POLL_INTERVAL. uses default 30", map[string]interface{}{
			"value": pollIntervalString,
			"error": err,
		})
		pollIntervalInMinutes = 30
	}

	logger.Info("Scheduler started", map[string]interface{}{
		"pollIntervalInMinutes": pollIntervalInMinutes,
	})

	executePollingCycle(db)

	ticker := time.NewTicker(time.Duration(pollIntervalInMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		executePollingCycle(db)
	}
}

func executePollingCycle(db *gorm.DB) {
	var watermarkRepo ports.PollWatermarkRepo = poll_watermarks_repo.NewGormRepo(db)

	taskTypes := []string{
		ProcessAmericanTournamentsTaskType,
		ProcessNotUrgentEuropeanTournamentsTaskType,
		ProcessUrgentEuropeanTournamentsTaskType,
	}

	now := time.Now().UTC()
	var wg sync.WaitGroup
	wg.Add(len(taskTypes))

	for _, taskType := range taskTypes {
		go func(taskType string) {
			defer wg.Done()

			task := matchExecutableAdapterByTaskType(taskType, db)
			if task == nil {
				logger.Error("Couldn't find executable adapter by task type", map[string]interface{}{
					"taskType": taskType,
				})
				return
			}

			startOfDay := toStartOfUTCDay(now)
			watermarkModel, err := watermarkRepo.FirstOrCreate(poll_watermarks.PollWatermarkModel{
				TaskType:             taskType,
				LastSuccessfulPollAt: startOfDay,
			})
			if err != nil {
				logger.Error("Couldn't read or create task watermark", map[string]interface{}{
					"taskType": taskType,
					"error":    err,
				})
				return
			}

			if err = task.Execute(startOfDay); err != nil {
				logger.Error("Error while processing tournament games", map[string]interface{}{
					"taskType": taskType,
					"error":    err,
				})
				return
			}

			watermarkModel.LastSuccessfulPollAt = now
			_, err = watermarkRepo.Update(watermarkModel)
			if err != nil {
				logger.Warn("Couldn't update task watermark", map[string]interface{}{
					"taskType": taskType,
					"error":    err,
				})
				return
			}

			logger.Info("Task executed successfully", map[string]interface{}{
				"taskType": taskType,
				"from":     startOfDay,
				"to":       now,
			})
		}(taskType)
	}

	wg.Wait()

	logger.Info("Polling cycle finished", map[string]interface{}{
		"finishedAt": now,
	})
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

func toStartOfUTCDay(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
