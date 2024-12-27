package mappers

import (
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/imp/results"
	"strconv"
)

func TeamToResult(team models.TeamGameResultModel, playersResults []results.PlayerFterResult) *results.TeamResults {
	return &results.TeamResults{
		Title:   team.Team.Alias + " - " + strconv.Itoa(team.TotalPoints),
		Players: playersResults,
	}
}
