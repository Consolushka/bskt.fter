package queries

import (
	"IMP/app/database"
	"IMP/app/internal/modules/games/models"
)

func SaveGameQuery(game models.GameModel) error {
	dbConnection := database.Connect()

	dbConnection.Create(&game)

	return dbConnection.Error
}
