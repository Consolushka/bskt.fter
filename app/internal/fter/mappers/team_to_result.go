package mappers

import (
	results2 "FTER/app/internal/fter/results"
	"FTER/app/internal/models"
	"strconv"
)

func TeamToResult(team models.TeamGameResultModel, playersResults []results2.PlayerFterResult) *results2.TeamResults {
	return &results2.TeamResults{
		Title:   team.Team.Alias + " - " + strconv.Itoa(team.TotalPoints),
		Players: playersResults,
	}
}
