package games_repo

import (
	"IMP/app/internal/core/games"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) GameExists(model games.GameModel) (bool, error) {
	var count int64
	tx := g.db.Model(&games.GameModel{}).Where(games.GameModel{
		TournamentId: model.TournamentId,
		Title:        model.Title,
		ScheduledAt:  model.ScheduledAt,
	}).Count(&count)

	return count > 0, tx.Error
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
