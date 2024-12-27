package commands

import (
	mappers2 "FTER/app/internal/modules/imp/mappers"
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/imp/queries"
	"FTER/app/internal/modules/imp/results"
	"FTER/app/internal/modules/statistics/enums"
)

// CalculateFullGame evaluates FTER for every player in each team
// and returns results.GameResult
func CalculateFullGame(game *models.GameModel) *results.GameResult {
	return mappers2.GameToResult(
		game,
		teamResults(game.Home, game.Away.TotalPoints, game.League),
		teamResults(game.Away, game.Home.TotalPoints, game.League),
	)
}

func teamResults(team models.TeamGameResultModel, oppPoints int, league enums.League) *results.TeamResults {
	playersFter := queries.FullTeamFter(team.Players, team.TotalPoints-oppPoints, league)

	return mappers2.TeamToResult(team, playersFter)
}
