package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"
	"strings"

	"gorm.io/gorm"
)

type GormRepo struct {
	db *gorm.DB
}

func (u GormRepo) ListTournamentsByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel
	conditionLeaguesAliases := "("

	for _, alias := range aliases {
		conditionLeaguesAliases += "'" + alias + "',"
	}
	conditionLeaguesAliases = strings.TrimRight(conditionLeaguesAliases, ",") + ")"

	err := u.db.Model(&tournaments.TournamentModel{}).Preload("League").Preload("ExternalIds").InnerJoins("left join leagues on tournaments.league_id = leagues.id").Find(&models, "leagues.alias in "+conditionLeaguesAliases).Error

	return models, err
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
