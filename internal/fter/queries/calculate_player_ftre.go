package queries

import (
	"FTER/internal/fter"
	"FTER/internal/fter/models"
	"FTER/internal/utils"
)

func PlayerFTRE(playerModel models.PlayerModel, finalDiff int, fullGameTime int) float64 {
	playedTime, err := utils.TimeToDecimal(playerModel.MinutesPlayed)
	if err != nil {
	}

	return fter.Calculate(playedTime, playerModel.PlsMin, finalDiff, fullGameTime)
}
