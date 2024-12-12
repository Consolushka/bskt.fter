package mappers

import (
	"FTER/internal/fter/models"
	"FTER/internal/fter/results"
	"strconv"
)

func TeamToResult(team models.TeamGameResultModel, playersResults []results.PlayerFterResult) *results.TeamResults {
	return &results.TeamResults{
		Title:   team.Team.Alias + " - " + strconv.Itoa(team.TotalPoints),
		Players: playersResults,
	}
}
