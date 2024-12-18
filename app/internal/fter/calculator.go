package fter

import (
	"FTER/internal/calculations"
)

func Calculate(playedTime float64, plsMin int, finalDiff int, fullGameTime int) float64 {
	onCourtEffPerMinute := float64(plsMin) / playedTime
	onCourtPerFullTime := onCourtEffPerMinute * float64(fullGameTime)

	rawValue := onCourtPerFullTime - float64(finalDiff)

	reliability := calculations.CalculateReliability(playedTime, float64(fullGameTime))

	return rawValue * reliability
}
