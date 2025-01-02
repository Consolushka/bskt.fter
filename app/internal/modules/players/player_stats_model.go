package players

import "time"

type PlayerGameStats struct {
	PlayerID      int       `json:"player_id" db:"player_id" gorm:"primaryKey;autoIncrement:false"`
	GameID        int       `json:"game_id" db:"game_id" gorm:"primaryKey;autoIncrement:false"`
	TeamID        int       `json:"team_id" db:"team_id"`
	PlsMin        int       `json:"pls_min" db:"pls_min"`
	PlayedSeconds int       `json:"played_min" db:"played_seconds"`
	IsBench       bool      `json:"is_bench" db:"is_bench"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
