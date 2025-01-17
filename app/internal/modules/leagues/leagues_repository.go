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

func (r *Repository) GetLeagueByAliasEn(aliasEn string) (*League, error) {
	var result League
	tx := r.db.Model(League{AliasEn: aliasEn}).First(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}
