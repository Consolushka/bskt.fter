package commands

import (
	"NBATrueEfficency/internal/fter/mappers"
	"NBATrueEfficency/internal/fter/models"
	"NBATrueEfficency/internal/fter/queries"
	"NBATrueEfficency/internal/fter/results"
)

// CalculateFullGame evaluates FTER for every player in each team
// and returns results.GameResult
func CalculateFullGame(game *models.GameModel) *results.GameResult {
	//todo: get duration from quarters
	gameDuration := game.League.FullGameDuration()

	return mappers.GameToResult(
		game,
		teamResults(game.Home, game.Away.TotalPoints, gameDuration),
		teamResults(game.Away, game.Home.TotalPoints, gameDuration),
	)
}

func teamResults(team models.TeamGameResultModel, oppPoints int, fullGameTime int) *results.TeamResults {
	playersFter := queries.FullTeamFter(team.Players, team.TotalPoints-oppPoints, fullGameTime)

	return mappers.TeamToResult(team, playersFter)
}
