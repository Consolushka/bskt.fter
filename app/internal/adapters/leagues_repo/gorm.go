package leagues_repo

import (
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/ports"
	"fmt"

	"gorm.io/gorm"
)

var _ ports.LeaguesRepo = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}

// FirstOrCreate находит лигу по паре (name, alias) — это её уникальный ключ —
// или создаёт новую. Идемпотентно: повторный вызов вернёт уже существующую.
func (g Gorm) FirstOrCreate(model leagues.LeagueModel) (leagues.LeagueModel, error) {
	tx := g.db.FirstOrCreate(&model, leagues.LeagueModel{
		Name:  model.Name,
		Alias: model.Alias,
	})
	if tx.Error != nil {
		return leagues.LeagueModel{}, fmt.Errorf("FirstOrCreate league %q returned error: %w", model.Name, tx.Error)
	}

	return model, nil
}
