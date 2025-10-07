package scheduler_repo

import (
	"IMP/app/internal/core/scheduler"
	"IMP/app/internal/ports"
	"time"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) TasksList() ([]scheduler.ScheduledTaskModel, error) {
	var result []scheduler.ScheduledTaskModel
	tx := g.db.Find(&result)

	return result, tx.Error
}

func (g Gorm) RescheduleTask(id uint, nextExecutionAt time.Time) (scheduler.ScheduledTaskModel, error) {
	var result scheduler.ScheduledTaskModel

	tx := g.db.First(&result, id)
	if tx.Error != nil {
		return result, tx.Error
	}

	result.NextExecutionAt = nextExecutionAt
	tx = g.db.Save(&result)

	return result, tx.Error
}

func NewGormRepo(db *gorm.DB) ports.SchedulerRepo {
	return Gorm{db: db}
}
