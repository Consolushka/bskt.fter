package models

import (
	"time"
)

type PlayerMetrics struct {
	PlayerID         int       `json:"player_id" db:"player_id" gorm:"primaryKey;type:uuid"`
	Player           *Player   `json:"player" gorm:"foreignKey:PlayerID"`
	AvgClearImp      float64   `json:"avg_clear_imp" db:"avg_clear_imp" gorm:"type:numeric(12,11)"`
	AvgPlayedSeconds float64   `json:"avg_played_seconds" db:"avg_played_time" gorm:"type:numeric(6,4)"`
	PlayedGamesCount int       `json:"played_games_count" db:"played_games_count"`
	FromBenchCount   int       `json:"from_bench_count" db:"from_bench_count"`
	FromStartCount   int       `json:"from_start_count" db:"from_start_count"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (PlayerMetrics) TableName() string {
	return "players_metrics"
}
