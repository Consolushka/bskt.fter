package resources

import (
	"IMP/app/internal/modules/players/domain/models"
	"IMP/app/internal/utils/time_utils"
)

type Player struct {
	FullName      string `json:"full_name"`
	MinutesPlayed string `json:"minutes_played"`
	Start         bool   `json:"start"`
	PlsMin        int    `json:"pls_min"`
}

func NewPlayerResource(playerGameStats models.PlayerGameStats) Player {
	formattedPlayedTime := time_utils.SecondsToFormat(playerGameStats.PlayedSeconds, time_utils.PlayedTimeFormat)

	return Player{
		FullName:      playerGameStats.Player.FullNameLocal,
		MinutesPlayed: formattedPlayedTime,
		Start:         !playerGameStats.IsBench,
		PlsMin:        playerGameStats.PlsMin,
	}
}
