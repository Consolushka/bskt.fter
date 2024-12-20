package queries

import (
	"FTER/app/internal/fter/mappers"
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
)

func FullTeamFter(players []models.PlayerModel, finalDiff int, fullGameTime int) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fter := PlayerFTRE(player, finalDiff, fullGameTime)
		tableData[i] = mappers.PlayerToResult(player, fter)
	}

	return tableData
}
