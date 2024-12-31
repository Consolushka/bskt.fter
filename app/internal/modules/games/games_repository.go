package games

import (
	"IMP/app/database"
)

func FirstOrCreate(game GameModel) (GameModel, error) {
	var result GameModel

	dbConnection := database.Connect()

	tx := dbConnection.
		Attrs(GameModel{
			PlayedMinutes: game.PlayedMinutes,
		}).
		FirstOrCreate(&result, GameModel{
			HomeTeamID:  game.HomeTeamID,
			AwayTeamID:  game.AwayTeamID,
			LeagueID:    game.LeagueID,
			ScheduledAt: game.ScheduledAt,
		})

	return result, tx.Error
}
