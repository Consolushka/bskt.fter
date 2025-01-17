package teams

import (
	"IMP/app/database"
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

func (r *Repository) FirstOrCreate(team Team) (Team, error) {
	var result Team

	tx := r.dbConnection.
		Attrs(
			Team{
				Name: team.Name,
			}).
		FirstOrCreate(&result,
			Team{
				Alias:    team.Alias,
				LeagueID: team.LeagueID,
			},
		)

	return result, tx.Error
}

func (r *Repository) FirstOrCreateGameStats(stats TeamGameStats) (TeamGameStats, error) {
	var result TeamGameStats

	tx := r.dbConnection.Attrs(
		TeamGameStats{
			Points: stats.Points,
		}).
		FirstOrCreate(
			&result,
			TeamGameStats{
				TeamId: stats.TeamId,
				GameId: stats.GameId,
			})

	return result, tx.Error
}
