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

func (r *Repository) FirstOrCreate(team TeamModel) (TeamModel, error) {
	var result TeamModel

	tx := r.dbConnection.
		Attrs(
			TeamModel{
				Name: team.Name,
			}).
		FirstOrCreate(&result,
			TeamModel{
				Alias:    team.Alias,
				LeagueID: team.LeagueID,
			},
		)

	return result, tx.Error
}
