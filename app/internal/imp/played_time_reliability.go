package imp

import (
	"math"
)

const lowerEdgeCoefficient = 0.3947

func CalculateReliability(minutesPlayed float64, timeBase TimeBasedImpCoefficient) float64 {
	timeBaseMinutes := float64(timeBase.Minutes())
	lowerEdge := lowerEdgeCoefficient * timeBaseMinutes

	if minutesPlayed < lowerEdge {
		return timeBase.InsufficientDistanceCoefficient() * math.Pow(minutesPlayed, timeBase.InsufficientDistancePower())
	}

	if minutesPlayed < timeBaseMinutes {
		return timeBase.SufficientDistanceOffset() + (math.Pow(minutesPlayed-lowerEdge, timeBase.SufficientDistancePower()))/(math.Pow(timeBaseMinutes-lowerEdge, timeBase.SufficientDistancePower()))*(1-timeBase.SufficientDistanceOffset())
	}

	res := 1 - (math.Pow(minutesPlayed-timeBaseMinutes, timeBase.OverSufficientDistanceUpperPower()))/(math.Pow(timeBaseMinutes-lowerEdge, timeBase.OverSufficientDistanceLowerPower()))
	return res
}
