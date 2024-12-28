package queries

import (
	"IMP/app/internal/modules/imp"
	mappers2 "IMP/app/internal/modules/imp/mappers"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/imp/results"
	"IMP/app/internal/modules/statistics/enums"
)

// CalculateFullGame evaluates IMP for every player in each team
// and returns results.GameResult
func CalculateFullGame(game *models.GameModel) *results.GameResult {
	return mappers2.GameToResult(
		game,
		teamResults(game.Home, game.Away.TotalPoints, game.League),
		teamResults(game.Away, game.Home.TotalPoints, game.League),
	)
}

func teamResults(team models.TeamGameResultModel, oppPoints int, league enums.League) *results.TeamResults {
	playersImps := imp.CalculateTeam(team.Players, team.TotalPoints-oppPoints, league)

	return mappers2.TeamToResult(team, playersImps)
}
