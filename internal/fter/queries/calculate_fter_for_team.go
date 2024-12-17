package queries

import (
	"FTER/internal/fter/mappers"
	"FTER/internal/fter/results"
	"FTER/internal/models"
)

func FullTeamFter(players []models.PlayerModel, finalDiff int, fullGameTime int) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fter := PlayerFTRE(player, finalDiff, fullGameTime)
		tableData[i] = mappers.PlayerToResult(player, fter)
	}

	return tableData
}
