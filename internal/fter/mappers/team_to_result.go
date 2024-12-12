package mappers

import (
	"NBATrueEfficency/internal/fter/models"
	"NBATrueEfficency/internal/fter/results"
	"strconv"
)

func TeamToResult(team models.TeamGameResultModel, playersResults []results.PlayerFterResult) *results.TeamResults {
	return &results.TeamResults{
		Title:   team.Team.Alias + " - " + strconv.Itoa(team.TotalPoints),
		Players: playersResults,
	}
}
