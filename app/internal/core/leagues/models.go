package leagues

import "time"

type LeagueModel struct {
	Id        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
