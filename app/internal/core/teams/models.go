package teams

import "time"

type TeamModel struct {
	Id        uint
	Name      string
	HomeTown  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GameTeamStatModel struct {
	Id        uint
	GameId    uint
	TeamId    uint
	Score     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
