package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/models"
)

func LeagueByAliasEn(aliasEn string) (models.League, error) {
	dbConn := database.Connect()

	var result models.League
	dbConn.Model(models.League{AliasEn: aliasEn}).First(&result)

	return result, nil
}
