package leagues

import (
	"IMP/app/database"
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

func (r *Repository) GetLeagueByAliasEn(aliasEn string) (League, error) {
	var result League
	r.db.Model(League{AliasEn: aliasEn}).First(&result)

	return result, nil
}
