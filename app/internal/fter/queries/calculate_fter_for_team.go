package queries

import (
	"FTER/app/internal/calculations"
	"FTER/app/internal/enums"
	enums2 "FTER/app/internal/fter/enums"
	"FTER/app/internal/fter/mappers"
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
	"FTER/app/internal/utils/time"
)

func FullTeamFter(players []models.PlayerModel, finalDiff int, league enums.League) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fullGameTime := league.FullGameTime()
		imp := PlayerFTRE(&player, finalDiff, fullGameTime)
		minutes, err := time.FromFormatToDecimal(player.MinutesPlayed)
		if err != nil {
		}

		bases := enums2.TimeBasesByLeague(league)
		impPers := make([]float64, len(bases))
		for i, timeBase := range bases {
			reliability := calculations.CalculateReliability(minutes, timeBase)
			pure := imp * float64(timeBase.Minutes())
			impPers[i] = pure * reliability
		}
		tableData[i] = mappers.PlayerToResult(player, bases, impPers)
	}

	return tableData
}
