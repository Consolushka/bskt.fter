package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"

	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

func (u GormRepo) ListTournamentsByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel
	if len(aliases) == 0 {
		return models, nil
	}

	err := u.db.Model(&tournaments.TournamentModel{}).
		Joins("JOIN leagues ON tournaments.league_id = leagues.id").
		Preload("League").
		Preload("Provider").
		Where("leagues.alias IN ?", aliases).
		Find(&models).Error

	return models, err
}

func NewGormRepo(db *gorm.DB) GormRepo {
	return GormRepo{db: db}
}

func (u GormRepo) ListActiveTournaments() ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel

	err := u.db.Preload("League").Preload("Provider").Find(&models).Error

	return models, err
}
