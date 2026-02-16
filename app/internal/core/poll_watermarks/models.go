package poll_watermarks

import "time"

type PollWatermarkModel struct {
	TournamentId         uint      `gorm:"column:tournament_id;primaryKey"`
	LastSuccessfulPollAt time.Time `gorm:"column:last_successful_poll_at"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (PollWatermarkModel) TableName() string {
	return "poll_watermarks"
}
