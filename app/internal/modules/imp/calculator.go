package imp

import (
	"FTER/app/internal/modules/calculations"
	enums2 "FTER/app/internal/modules/imp/enums"
	"FTER/app/internal/modules/imp/mappers"
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/imp/results"
	"FTER/app/internal/modules/statistics/enums"
)

// calculatePlayerImpPerMinute calculates IMPPerMinute for player
func calculatePlayerImpPerMinute(playedTime float64, plsMin int, finalDiff int, fullGameTime int) float64 {
	playerImpPerMinute := float64(plsMin) / playedTime
	fullGameImpPerMinute := float64(finalDiff) / float64(fullGameTime)

	rawValue := playerImpPerMinute - fullGameImpPerMinute

	return rawValue
}

// CalculateTeam calculates IMP for every player in the team for league time bases
func CalculateTeam(players []models.PlayerModel, finalDiff int, league enums.League) []results.PlayerFterResult {
	tableData := make([]results.PlayerFterResult, len(players))
	for i, player := range players {
		fullGameTime := league.FullGameTimeMinutes()
		playedMinutes := float64(player.SecondsPlayed) / 60

		impPerMinute := calculatePlayerImpPerMinute(playedMinutes, player.PlsMin, finalDiff, fullGameTime)

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
