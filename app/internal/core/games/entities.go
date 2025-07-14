package games

import "IMP/app/internal/core/teams"

type GameStatEntity struct {
	GameModel    GameModel
	HomeTeamStat teams.GameTeamStatModel
	AwayTeamStat teams.GameTeamStatModel
}
