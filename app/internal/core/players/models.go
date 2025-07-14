package players

import "time"

type PlayerModel struct {
	Id        int
	FullName  string
	BirthDate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GameTeamPlayerStatModel struct {
	Id            uint
	GameTeamId    uint
	PlayerId      uint
	PlayedSeconds int
	PlsMin        int8
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
