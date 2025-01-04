package players

import (
	"IMP/app/database"
)

func FirstOrCreate(player Player) (Player, error) {
	var result Player
	dbConnection := database.Connect()

	tx := dbConnection.
		FirstOrCreate(
			&result,
			Player{
				FullName:  player.FullName,
				BirthDate: player.BirthDate,
				DraftYear: player.DraftYear,
			})

	return result, tx.Error
}

func FirstOrCreateGameStat(stats PlayerGameStats) error {
	dbConnection := database.Connect()

	tx := dbConnection.FirstOrCreate(
		&PlayerGameStats{},
		PlayerGameStats{
			PlayerID: stats.PlayerID,
			GameID:   stats.GameID,
			TeamID:   stats.TeamID,
		}).Attrs(
		PlayerGameStats{
			PlayedSeconds: stats.PlayedSeconds,
			PlsMin:        stats.PlsMin,
			IsBench:       stats.IsBench,
		})

	return tx.Error
}
