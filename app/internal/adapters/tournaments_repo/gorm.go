package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"
	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) GormRepo {
	return GormRepo{db: db}
}

func (u GormRepo) ListActiveTournaments() ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel

	err := u.db.Preload("League").Find(&models).Error

	return models, err
}
