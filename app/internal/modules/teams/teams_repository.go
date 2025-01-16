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
