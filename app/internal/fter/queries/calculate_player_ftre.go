package queries

import (
	"FTER/app/internal/fter"
	"FTER/app/internal/models"
	"FTER/app/internal/utils/time"
)

func PlayerFTRE(playerModel *models.PlayerModel, finalDiff int, fullGameTime int) float64 {
	playedTime, err := time.FromFormatToDecimal(playerModel.MinutesPlayed)
	if err != nil {
	}

	return fter.Calculate(playedTime, playerModel.PlsMin, finalDiff, fullGameTime)
}
