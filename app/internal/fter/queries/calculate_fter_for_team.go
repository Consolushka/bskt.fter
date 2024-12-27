package queries

import (
	"FTER/app/internal/calculations"
	"FTER/app/internal/enums"
	"FTER/app/internal/fter"
	enums2 "FTER/app/internal/fter/enums"
	"FTER/app/internal/fter/mappers"
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
)

func FullTeamFter(players []models.PlayerModel, finalDiff int, league enums.League) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fullGameTime := league.FullGameTimeMinutes()
		playedMinutes := float64(player.SecondsPlayed) / 60

		impPerMinute := fter.Calculate(playedMinutes, player.PlsMin, finalDiff, fullGameTime)

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
