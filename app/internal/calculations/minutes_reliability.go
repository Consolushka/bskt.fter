package calculations

import "math"

const lowerEdgeCoef = 0.3125
const sampleGameTimeCoef = 0.79

func CalculateReliability(minutesPlayed float64, fullGameTime float64) float64 {
	lowerEdge := lowerEdgeCoef * fullGameTime
	sampleGameTime := sampleGameTimeCoef * fullGameTime

	if minutesPlayed < lowerEdge {
		return 0.00012 * math.Pow(minutesPlayed, 3)
	}

	if minutesPlayed < sampleGameTime {
		return 0.405 + (math.Pow(minutesPlayed-lowerEdge, 0.6))/(math.Pow(sampleGameTime-lowerEdge, 0.6))*(1-0.405)

	}

	return 1 - (math.Pow(minutesPlayed-sampleGameTime, 3))/(math.Pow(sampleGameTime-lowerEdge, 3))
}
