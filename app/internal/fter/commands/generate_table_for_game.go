package commands

import (
	"FTER/app/internal/enums"
	"FTER/app/internal/fter/mappers"
	"FTER/app/internal/fter/queries"
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
)

// CalculateFullGame evaluates FTER for every player in each team
// and returns results.GameResult
func CalculateFullGame(game *models.GameModel) *results.GameResult {
	return mappers.GameToResult(
		game,
		teamResults(game.Home, game.Away.TotalPoints, game.League),
		teamResults(game.Away, game.Home.TotalPoints, game.League),
	)
}

func teamResults(team models.TeamGameResultModel, oppPoints int, league enums.League) *results.TeamResults {
	playersFter := queries.FullTeamFter(team.Players, team.TotalPoints-oppPoints, league)

	return mappers.TeamToResult(team, playersFter)
}
