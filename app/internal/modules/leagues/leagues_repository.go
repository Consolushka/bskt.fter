package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/models"
)

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r Repository) LeagueByAliasEn(aliasEn string) (models.League, error) {
	dbConn := database.Connect()

	var result models.League
	dbConn.Model(models.League{AliasEn: aliasEn}).First(&result)

	return result, nil
}
