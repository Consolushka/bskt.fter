package teams

import (
	"time"
)

type Team struct {
	ID        int       `json:"id" db:"id"`
	Alias     string    `json:"alias" db:"alias"`
	LeagueID  int       `json:"league_id" db:"league_id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Team) TableName() string {
	return "teams"
}
