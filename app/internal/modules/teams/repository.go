package teams

import (
	"IMP/app/database"
	"IMP/app/internal/modules/teams/models"
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

func (r *Repository) FirstOrCreate(team models.Team) (models.Team, error) {
	var result models.Team

	tx := r.dbConnection.
		Attrs(
			models.Team{
				Name:       team.Name,
				OfficialId: team.OfficialId,
			}).
		FirstOrCreate(&result,
			models.Team{
				Alias:    team.Alias,
				LeagueID: team.LeagueID,
			},
		)

	return result, tx.Error
}

func (r *Repository) FirstOrCreateGameStats(stats models.TeamGameStats) (models.TeamGameStats, error) {
	var result models.TeamGameStats

	tx := r.dbConnection.Attrs(
		models.TeamGameStats{
			Points: stats.Points,
		}).
		FirstOrCreate(
			&result,
			models.TeamGameStats{
				TeamId: stats.TeamId,
				GameId: stats.GameId,
			})

	return result, tx.Error
}
