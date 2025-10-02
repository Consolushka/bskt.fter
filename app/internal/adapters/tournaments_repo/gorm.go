package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"

	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

func (u GormRepo) FindTournamentExternalId(tournamentId uint, providerName string) (tournaments.TournamentExternalIdModel, error) {
	var model tournaments.TournamentExternalIdModel

	err := u.db.Preload("Tournament").Find(&model, tournaments.TournamentExternalIdModel{TournamentId: tournamentId, ProviderName: providerName}).Error

	return model, err
}

func NewGormRepo(db *gorm.DB) GormRepo {
	return GormRepo{db: db}
}

func (u GormRepo) ListActiveTournaments() ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel

	err := u.db.Preload("League").Preload("ExternalIds").Find(&models).Error

	return models, err
}
