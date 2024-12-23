package queries

import (
	"FTER/app/internal/fter"
	"FTER/app/internal/models"
	"FTER/app/internal/utils"
)

func PlayerFTRE(playerModel *models.PlayerModel, finalDiff int, fullGameTime int) float64 {
	playedTime, err := utils.TimeToDecimal(playerModel.MinutesPlayed)
	if err != nil {
	}

	return fter.Calculate(playedTime, playerModel.PlsMin, finalDiff, fullGameTime)
}
