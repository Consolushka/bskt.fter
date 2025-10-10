package scheduler

import (
	"time"

	"gorm.io/gorm"
)

type ScheduledTaskModel struct {
	Id              uint           `gorm:"column:id"`
	Type            string         `gorm:"column:type"`
	LastExecutedAt  time.Time      `gorm:"column:last_executed_at"`
	NextExecutionAt time.Time      `gorm:"column:next_execution_at"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ScheduledTaskModel) TableName() string {
	return "scheduled_tasks"
}
