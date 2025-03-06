package domain

import (
	"IMP/app/database"
	teamsModels "IMP/app/internal/modules/teams/domain/models"
	"gorm.io/gorm"
)

type Repository struct {
	dbConnection *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{
		dbConnection: database.GetDB(),
	}
}

func (r *Repository) FirstOrCreate(team teamsModels.Team) (teamsModels.Team, error) {
	var result teamsModels.Team

	tx := r.dbConnection.
		Attrs(
			teamsModels.Team{
				Name:       team.Name,
				OfficialId: team.OfficialId,
			}).
		FirstOrCreate(&result,
			teamsModels.Team{
				Alias:    team.Alias,
				LeagueID: team.LeagueID,
			},
		)

	return result, tx.Error
}

func (r *Repository) FirstOrCreateTeamGameStats(stats teamsModels.TeamGameStats) (teamsModels.TeamGameStats, error) {
	var result teamsModels.TeamGameStats

	tx := r.dbConnection.Attrs(
		teamsModels.TeamGameStats{
			Points: stats.Points,
		}).
		FirstOrCreate(
			&result,
			teamsModels.TeamGameStats{
				TeamId: stats.TeamId,
				GameId: stats.GameId,
			})

	return result, tx.Error
}

func (r *Repository) List() ([]teamsModels.Team, error) {
	var teams []teamsModels.Team

	tx := r.dbConnection.Model(&teamsModels.Team{}).
		Preload("League").
		Find(&teams)

	return teams, tx.Error
}

func (r *Repository) FirstById(id int) (teamsModels.Team, error) {
	var team teamsModels.Team

	tx := r.dbConnection.Model(&teamsModels.Team{}).
		Preload("League").
		Where("id = ?", id).
		First(&team)

	return team, tx.Error
}

func (r *Repository) TeamGameIdListByTeamId(teamId int) ([]int, error) {
	var gamesIds []int

	tx := r.dbConnection.Model(&teamsModels.TeamGameStats{}).
		Where("team_id = ?", teamId).
		Select("game_id").
		Find(&gamesIds)

	return gamesIds, tx.Error
}

func (r *Repository) ListByLeague(leagueId int) ([]teamsModels.Team, error) {
	var teams []teamsModels.Team

	tx := r.dbConnection.
		Where("league_id = ?", leagueId).
		Find(&teams)

	return teams, tx.Error
}
