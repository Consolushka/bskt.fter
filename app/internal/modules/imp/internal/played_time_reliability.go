package calculations

import (
	"IMP/app/internal/modules/imp/internal/domain"
	"math"
)

const lowerEdgeCoef = 0.3947

func CalculateReliability(minutesPlayed float64, timeBase domain.TimeBasedImpCoefficient) float64 {
	timeBaseMinutes := float64(timeBase.Minutes())
	lowerEdge := lowerEdgeCoef * timeBaseMinutes

	if minutesPlayed < lowerEdge {
		return timeBase.InsufficientDistanceCoef() * math.Pow(minutesPlayed, timeBase.InsufficientDistancePower())
	}

	if minutesPlayed < timeBaseMinutes {
		return timeBase.SufficientDistanceOffset() + (math.Pow(minutesPlayed-lowerEdge, timeBase.SufficientDistancePower()))/(math.Pow(timeBaseMinutes-lowerEdge, timeBase.SufficientDistancePower()))*(1-timeBase.SufficientDistanceOffset())
	}

	res := 1 - (math.Pow(minutesPlayed-timeBaseMinutes, timeBase.OverSufficientDistanceUpperPower()))/(math.Pow(timeBaseMinutes-lowerEdge, timeBase.OverSufficientDistanceLowerPower()))
	return res
}
