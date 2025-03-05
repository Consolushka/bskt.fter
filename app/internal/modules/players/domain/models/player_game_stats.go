package models

import (
	"time"
)

type PlayerGameStats struct {
	PlayerID      int       `json:"player_id" db:"player_id" gorm:"primaryKey;autoIncrement:false"`
	Player        *Player   `json:"player" gorm:"foreignKey:PlayerID"`
	TeamGameId    int       `json:"team_game_id" db:"team_game_id" gorm:"index"`
	PlsMin        int       `json:"pls_min" db:"pls_min"`
	PlayedSeconds int       `json:"played_min" db:"played_seconds"`
	IsBench       bool      `json:"is_bench" db:"is_bench"`
	IMPClean      float64   `json:"imp_clean" db:"imp_clean"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
