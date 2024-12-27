package queries

import (
	"FTER/app/internal/modules/calculations"
	"FTER/app/internal/modules/imp"
	enums2 "FTER/app/internal/modules/imp/enums"
	"FTER/app/internal/modules/imp/mappers"
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/imp/results"
	"FTER/app/internal/modules/statistics/enums"
)

func FullTeamFter(players []models.PlayerModel, finalDiff int, league enums.League) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fullGameTime := league.FullGameTimeMinutes()
		playedMinutes := float64(player.SecondsPlayed) / 60

		impPerMinute := imp.Calculate(playedMinutes, player.PlsMin, finalDiff, fullGameTime)

		bases := enums2.TimeBasesByLeague(league)
		impPers := make([]float64, len(bases))
		for i, timeBase := range bases {
			reliability := calculations.CalculateReliability(playedMinutes, timeBase)
			pure := impPerMinute * float64(timeBase.Minutes())
			impPers[i] = pure * reliability
		}
		tableData[i] = mappers.PlayerToResult(player, bases, impPers)
	}

	return tableData
}
