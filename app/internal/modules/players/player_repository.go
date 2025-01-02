package players

import (
	"IMP/app/database"
)

func FirstOrCreate(player Player) (Player, error) {
	var result Player
	dbConnection := database.Connect()

	tx := dbConnection.
		FirstOrCreate(&result, Player{FullName: player.FullName, BirthDate: player.BirthDate})

	return result, tx.Error
}

func CreateStatisticInGame(stats PlayerGameStats) error {
	dbConnection := database.Connect()

	tx := dbConnection.Create(&stats)

	return tx.Error
}
