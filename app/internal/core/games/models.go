package games

import "time"

type GameModel struct {
	Id           uint
	TournamentId uint
	ScheduledAt  time.Time
	Title        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
