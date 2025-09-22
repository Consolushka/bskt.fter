package games_repo

import (
	"IMP/app/internal/core/games"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) FindOrCreateGame(model games.GameModel) (games.GameModel, error) {
	tx := g.db.FirstOrCreate(&model, games.GameModel{
		TournamentId: model.TournamentId,
		Title:        model.Title,
		ScheduledAt:  model.ScheduledAt,
	})

	return model, tx.Error
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
