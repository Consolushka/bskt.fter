package models

import (
	"NBATrueEfficency/internal/enums"
)

type TeamGameResultModel struct {
	Team        TeamModel
	TotalPoints int
	Players     []PlayerModel
}

type GameModel struct {
	Id        string
	Scheduled string
	Home      TeamGameResultModel
	Away      TeamGameResultModel
	League    enums.League
}
