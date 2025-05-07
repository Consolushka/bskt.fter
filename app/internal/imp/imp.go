package imp

import (
	"IMP/app/internal/domain"
)

type ImpPERs string

const (
	Clean    ImpPERs = "Clean"
	Bench    ImpPERs = "Bench"
	Start    ImpPERs = "Start"
	FullGame ImpPERs = "FullGame"
)

func (ip ImpPERs) Order() int {
	switch ip {
	case Clean:
		return 0
	case Bench:
		return 1
	case Start:
		return 2
	case FullGame:
		return 3
	}

	return -1
}

func (ip ImpPERs) ToString() string {
	return string(ip)
}

func EvaluateClean(playedSeconds int, plsMin int, finalDifference int, fullGameTimeInMinutes int) float64 {
	if playedSeconds == 0 {
		return 0
	}

	playerImpPerMinute := float64(plsMin) / (float64(playedSeconds) / 60)
	fullGameImpPerMinute := float64(finalDifference) / float64(fullGameTimeInMinutes)

	rawValue := playerImpPerMinute - fullGameImpPerMinute

	return rawValue
}

func EvaluatePer(playedSeconds int, plsMin *int, finalDifferential *int, fullGameTime *int, impPer ImpPERs, league *domain.League, cleanImpPointer *float64) float64 {
	var cleanImp float64

	if cleanImpPointer == nil {
		cleanImp = EvaluateClean(playedSeconds, *plsMin, *finalDifferential, *fullGameTime)
	} else {
		cleanImp = *cleanImpPointer
	}

	timeBase := TimeBasesByLeagueAndPers(league, impPer)

	reliability := CalculateReliability(float64(playedSeconds)/60, timeBase)
	pure := cleanImp * float64(timeBase.Minutes())

	return pure * reliability
}
