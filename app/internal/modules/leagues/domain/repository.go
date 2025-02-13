package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/domain/models"
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

func (r *Repository) GetLeagueByAliasEn(aliasEn string) (*models.League, error) {
	var result models.League
	tx := r.db.Model(models.League{AliasEn: aliasEn}).First(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}
