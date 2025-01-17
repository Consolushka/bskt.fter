package teams

import "time"

type TeamGameStats struct {
	TeamId    int       `db:"team_id"`
	GameId    int       `db:"game_id"`
	Points    int       `db:"points"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
