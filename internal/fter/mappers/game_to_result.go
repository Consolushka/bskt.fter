package mappers

import (
	"NBATrueEfficency/internal/fter/models"
	"NBATrueEfficency/internal/fter/results"
)

func GameToResult(game *models.GameModel, home *results.TeamResults, away *results.TeamResults) *results.GameResult {
	return &results.GameResult{
		GameId:   game.Id,
		Title:    game.Home.Team.Alias + " - " + game.Away.Team.Alias + ". " + game.Scheduled,
		Schedule: game.Scheduled,
		Home:     home,
		Away:     away,
	}
}
