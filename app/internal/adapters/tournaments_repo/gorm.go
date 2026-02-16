package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) ListByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel
	if len(aliases) == 0 {
		return models, nil
	}

	err := g.db.Model(&tournaments.TournamentModel{}).
		Joins("JOIN leagues ON tournaments.league_id = leagues.id").
		Preload("League").
		Preload("Provider").
		Where("leagues.alias IN ?", aliases).
		Find(&models).Error

	return models, err
}

func (g Gorm) Get(id uint) (tournaments.TournamentModel, error) {
	var model tournaments.TournamentModel
	err := g.db.Preload("League").Preload("Provider").First(&model, id).Error
	return model, err
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}

func (g Gorm) ListActive() ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel

	err := g.db.Preload("League").Preload("Provider").Find(&models).Error

	return models, err
}
