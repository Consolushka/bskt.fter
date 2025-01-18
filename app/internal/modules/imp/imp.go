package imp

import (
	"IMP/app/internal/modules/calculations"
	enums2 "IMP/app/internal/modules/imp/enums"
	"IMP/app/internal/modules/imp/mappers"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/imp/results"
	"IMP/app/internal/modules/statistics/enums"
)

// CalculatePlayerImpPerMinute calculates IMPPerMinute for player
func CalculatePlayerImpPerMinute(playedTime float64, plsMin int, finalDiff int, fullGameTime int) float64 {
	if playedTime == 0 {
		return 0
	}

	playerImpPerMinute := float64(plsMin) / playedTime
	fullGameImpPerMinute := float64(finalDiff) / float64(fullGameTime)

	rawValue := playerImpPerMinute - fullGameImpPerMinute

	return rawValue
}

// CalculateTeam calculates IMP for every player in the team for league time bases
func CalculateTeam(players []models.PlayerModel, finalDiff int, league enums.League) []results.PlayerImpResult {
	tableData := make([]results.PlayerImpResult, len(players))
	for i, player := range players {
		fullGameTime := league.FullGameTimeMinutes()
		playedMinutes := float64(player.SecondsPlayed) / 60

		impPerMinute := CalculatePlayerImpPerMinute(playedMinutes, player.PlsMin, finalDiff, fullGameTime)

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
