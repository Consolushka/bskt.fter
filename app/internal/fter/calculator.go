package fter

func Calculate(playedTime float64, plsMin int, finalDiff int, fullGameTime int) float64 {
	//onCourtEffPerMinute := float64(plsMin) / playedTime
	//sampleGameTime := float64(38)
	//onCourtPerFullTime := onCourtEffPerMinute * sampleGameTime

	playerImpPerMinute := float64(plsMin) / playedTime
	fullGameImpPerMinute := float64(finalDiff) / float64(fullGameTime)

	rawValue := playerImpPerMinute - fullGameImpPerMinute

	//reliability := calculations.CalculateReliability(playedTime, float64(fullGameTime), sampleGameTime)

	//return rawValue * reliability
	return rawValue
}
