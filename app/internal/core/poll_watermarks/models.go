package poll_watermarks

import "time"

type PollWatermarkModel struct {
	TaskType             string    `gorm:"column:task_type"`
	LastSuccessfulPollAt time.Time `gorm:"column:last_successful_poll_at"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (PollWatermarkModel) TableName() string {
	return "poll_watermarks"
}
