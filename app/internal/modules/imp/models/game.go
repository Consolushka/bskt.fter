package models

import (
	"IMP/app/internal/modules/statistics/enums"
)

type TeamGameResultModel struct {
	Team        TeamModel
	TotalPoints int
	Players     []PlayerModel
}

type GameModel struct {
	Id           string
	Scheduled    string
	Home         TeamGameResultModel
	Away         TeamGameResultModel
	FullGameTime int
	League       enums.League
}
