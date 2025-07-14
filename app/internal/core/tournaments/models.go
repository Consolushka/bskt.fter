package tournaments

import "time"

type TournamentModel struct {
	Id        uint
	LeagueId  uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
