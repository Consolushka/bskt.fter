package teams_repo

import (
	"IMP/app/internal/core/teams"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) FirstOrCreate(model teams.TeamModel) (teams.TeamModel, error) {
	tx := g.db.FirstOrCreate(&model, teams.TeamModel{
		Name:     model.Name,
		HomeTown: model.HomeTown,
	})

	return model, tx.Error
}

func (g Gorm) FirstOrCreateStats(model teams.GameTeamStatModel) (teams.GameTeamStatModel, error) {
	tx := g.db.FirstOrCreate(&model, teams.GameTeamStatModel{
		GameId: model.GameId,
		TeamId: model.TeamId,
	})

	return model, tx.Error
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
