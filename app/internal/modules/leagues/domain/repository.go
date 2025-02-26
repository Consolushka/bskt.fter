package domain

import (
	"IMP/app/database"
	"IMP/app/internal/modules/leagues/domain/models"
	"gorm.io/gorm"
	"strings"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

func (r *Repository) List() ([]models.League, error) {
	var result []models.League

	tx := r.db.Find(&result)

	return result, tx.Error
}

func (r *Repository) FirstByAliasEn(aliasEn string) (*models.League, error) {
	var result models.League
	tx := r.db.First(&result, models.League{AliasEn: strings.ToUpper(aliasEn)})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}

func (r *Repository) FirstById(id int) (*models.League, error) {
	var leagueModel models.League

	tx := r.db.First(&leagueModel, models.League{ID: id})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &leagueModel, nil
}
