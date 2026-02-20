package tournament_poll_logs

import (
	"time"
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type TournamentPollLogModel struct {
	Id              uint64     `gorm:"column:id;primaryKey"`
	TournamentId    uint       `gorm:"column:tournament_id"`
	PollStartAt     time.Time  `gorm:"column:poll_start_at"`
	PollEndAt       *time.Time `gorm:"column:poll_end_at"`
	IntervalStart   time.Time  `gorm:"column:interval_start"`
	IntervalEnd     time.Time  `gorm:"column:interval_end"`
	SavedGamesCount int        `gorm:"column:saved_games_count"`
	Status          string     `gorm:"column:status"`
	ErrorMessage    *string    `gorm:"column:error_message"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
}

func (TournamentPollLogModel) TableName() string {
	return "tournament_poll_logs"
}
