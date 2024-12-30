package models

import "time"

type GameModel struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	LeagueID      int       `json:"league_id" gorm:"not null"`
	PlayedMinutes int       `json:"played_minutes" gorm:"not null"`
	ScheduledAt   time.Time `json:"scheduled_at" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (GameModel) TableName() string {
	return "games"
}
