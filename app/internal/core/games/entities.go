package games

import "IMP/app/internal/core/teams"

type GameStatEntity struct {
	GameModel          GameModel
	ExternalGameId     string
	HomeTeamExternalId int
	AwayTeamExternalId int
	HomeTeamStat       teams.TeamStatEntity
	AwayTeamStat       teams.TeamStatEntity
}
