package imp

import (
	"IMP/app/internal/modules/imp/domain/enums"
	calculations "IMP/app/internal/modules/imp/internal"
	enums3 "IMP/app/internal/modules/imp/internal/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
)

func EvaluateClean(playedSeconds int, plsMin int, finalDifference int, fullGameTime int) float64 {
	if playedSeconds == 0 {
		return 0
	}

	playerImpPerMinute := float64(plsMin) / (float64(playedSeconds) / 60)
	fullGameImpPerMinute := float64(finalDifference) / float64(fullGameTime)

	rawValue := playerImpPerMinute - fullGameImpPerMinute

	return rawValue
}

func EvaluatePer(playedSeconds int, plsMin int, finalDifferential int, fullGameTime int, impPer enums.ImpPERs, league leaguesModels.League, cleanImpPointer *float64) float64 {
	var cleanImp float64

	if cleanImpPointer == nil {
		cleanImp = EvaluateClean(playedSeconds, plsMin, finalDifferential, fullGameTime)
	} else {
		cleanImp = *cleanImpPointer
	}

	timeBase := enums3.TimeBasesByLeagueAndPers(league, impPer)

	reliability := calculations.CalculateReliability(float64(playedSeconds)/60, timeBase)
	pure := cleanImp * float64(timeBase.Minutes())

	return pure * reliability
}
