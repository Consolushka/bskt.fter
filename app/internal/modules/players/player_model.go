package players

import "time"

type Player struct {
	ID             int        `json:"id" db:"id"`
	FullName       string     `json:"full_name" db:"full_name"`
	BirthDate      *time.Time `json:"birth_date" db:"birth_date"`
	LeaguePlayerID int        `json:"league_player_id" db:"league_player_id"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
