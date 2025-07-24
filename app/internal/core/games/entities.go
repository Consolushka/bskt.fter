package games

import "IMP/app/internal/core/teams"

type GameStatEntity struct {
	GameModel    GameModel
	HomeTeamStat teams.TeamStatEntity
	AwayTeamStat teams.TeamStatEntity
}
