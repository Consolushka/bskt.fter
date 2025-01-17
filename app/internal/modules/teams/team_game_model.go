package teams

import (
	"IMP/app/internal/modules/players"
	"time"
)

type TeamGameStats struct {
	Id              int                       `json:"id" db:"id"`
	TeamId          int                       `json:"team_id" db:"team_id"`
	Team            Team                      `json:"team" gorm:"foreignKey:TeamId"`
	GameId          int                       `json:"game_id" db:"game_id"`
	Points          int                       `json:"points" db:"points"`
	PlayerGameStats []players.PlayerGameStats `json:"players" gorm:"foreignKey:team_game_id"`
	CreatedAt       time.Time                 `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time                 `json:"updated_at" db:"updated_at"`
}
