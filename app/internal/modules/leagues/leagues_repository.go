package leagues

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

func (r *Repository) GetLeagueByAliasEn(aliasEn string) (models.League, error) {
	var result models.League
	r.db.Model(models.League{AliasEn: aliasEn}).First(&result)

	return result, nil
}
