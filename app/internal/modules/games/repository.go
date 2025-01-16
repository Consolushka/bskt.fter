package games

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

func (r *Repository) First(id int) (*GameModel, error) {
	var result GameModel

	tx := r.dbConnection.
		First(&result, GameModel{ID: id}).
		Preload("League")

	return &result, tx.Error
}

func (r *Repository) FirstOrCreate(game GameModel) (GameModel, error) {
	var result GameModel

	tx := r.dbConnection.
		Attrs(GameModel{
			PlayedMinutes: game.PlayedMinutes,
		}).
		FirstOrCreate(&result, GameModel{
			HomeTeamID:  game.HomeTeamID,
			AwayTeamID:  game.AwayTeamID,
			LeagueID:    game.LeagueID,
			ScheduledAt: game.ScheduledAt,
		})

	return result, tx.Error
}
